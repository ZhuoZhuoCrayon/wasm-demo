/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/binary"
	"errors"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"golang.org/x/net/http2/hpack"
	// 将 nottinygc 作为 TinyGo 编译 WASI 的一个替代内存分配器，默认的内存分配器在数据量大的场景会有性能问题
	_ "github.com/wasilibs/nottinygc"
)

type FrameType uint8

type Flags uint8

func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

const (
	GRPC                uint8     = 41
	FrameHeaderLen      int       = 9
	FrameTypeMax        FrameType = 0x9
	FrameHeaders        FrameType = 0x1
	FrameContinuation   FrameType = 0x9
	FlagHeadersPadded   Flags     = 0x8
	FlagHeadersPriority Flags     = 0x20
)

var (
	ParsePayloadError = errors.New("failed to parse Payload")
	ExtractError      = errors.New("failed to extract headers")
)

var ExpectedHeaderFields = map[string]bool{
	"openid": true,
}

type FrameHeader struct {
	Length uint32
	Type   FrameType
	Flags  Flags
}

func readByte(p []byte) (remain []byte, b byte, err error) {
	if len(p) == 0 {
		return nil, 0, ParsePayloadError
	}
	return p[1:], p[0], nil
}

func readUint32(p []byte) (remain []byte, v uint32, err error) {
	if len(p) < 4 {
		return nil, 0, ParsePayloadError
	}
	return p[4:], binary.BigEndian.Uint32(p[:4]), nil
}

func readFrameHeader(p []byte) ([]byte, FrameHeader, error) {
	if len(p) < FrameHeaderLen {
		return nil, FrameHeader{}, ParsePayloadError
	}

	fh := FrameHeader{
		Type: FrameType(p[3]),
	}
	if fh.Type > FrameTypeMax {
		return nil, fh, ParsePayloadError
	}

	fh.Length = uint32(p[0])<<16 | uint32(p[1])<<8 | uint32(p[2])
	fh.Flags = Flags(p[4])

	return p[FrameHeaderLen:], fh, nil
}

func frameOverflow(p []byte, fh FrameHeader) bool {
	if int(fh.Length) > len(p) {
		return true
	}
	return false
}

func readHeaderBlockFragment(p []byte, fh FrameHeader) (remain []byte, bf []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	var padLength uint8
	if fh.Flags.Has(FlagHeadersPadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, nil, err
		}
	}
	if fh.Flags.Has(FlagHeadersPriority) {
		// read E + Stream Dependency
		if p, _, err = readUint32(p); err != nil {
			return nil, nil, err
		}
		// read Weight
		if p, _, err = readByte(p); err != nil {
			return nil, nil, err
		}
	}

	if len(p)-int(padLength) < 0 {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:len(p)-int(padLength)], nil
}

func readFramePayload(p []byte, fh FrameHeader) (remain []byte, fp []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:fh.Length], nil
}

// parse 解析 HTTP2 payload
func parse(payload []byte) (headers map[string]string, err error) {
	sdk.Warn("start to parse Payload: %x", payload)
	headers = make(map[string]string, len(ExpectedHeaderFields))
	hd := hpack.NewDecoder(4096, func(hf hpack.HeaderField) {
		if _, ok := ExpectedHeaderFields[hf.Name]; ok {
			headers[hf.Name] = hf.Value
		}
	})

	p := payload
	var fh FrameHeader

	// 输入 wasm 插件的 Payload 可能存在数据截断的情况，使用 break 结束循环而不是 return nil, err 来尽可能匹配前面符合规则的帧
	for len(p) > FrameHeaderLen {
		if p, fh, err = readFrameHeader(p); err != nil {
			break
		}
		sdk.Info("[parse] FrameType -> %#x", fh.Type)

		if fh.Type == FrameHeaders {
			var headerBlockFragment []byte
			if p, headerBlockFragment, err = readHeaderBlockFragment(p, fh); err != nil {
				break
			}
			if _, err = hd.Write(headerBlockFragment); err != nil {
				break
			}
		} else if fh.Type == FrameContinuation {
			var headerBlockFragment []byte
			if p, headerBlockFragment, err = readFramePayload(p, fh); err != nil {
				break
			}
			if _, err = hd.Write(headerBlockFragment); err != nil {
				break
			}
		} else {
			if p, _, err = readFramePayload(p, fh); err != nil {
				break
			}
		}
	}
	return headers, nil
}

type Extractor struct {
	Payload []byte
}

func (e *Extractor) do() ([]sdk.KeyVal, *sdk.Trace, error) {
	headers, err := parse(e.Payload)
	if err != nil {
		return nil, nil, ExtractError
	}
	sdk.Info("extract headers: %v", headers)
	if len(headers) == 0 {
		return nil, nil, ExtractError
	}

	attrs := make([]sdk.KeyVal, 0, len(headers))
	for k, v := range headers {
		attrs = append(attrs, sdk.KeyVal{Key: k, Val: v})
	}

	trace := &sdk.Trace{
		TraceID: attrs[0].Val,
	}
	return attrs, trace, nil
}

type httpHook struct {
}

func (p httpHook) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
		sdk.HOOK_POINT_HTTP_RESP,
	}
}

func (p httpHook) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != GRPC || baseCtx.DstPort != 9000 {
		return sdk.ActionNext()
	}
	sdk.Warn("[Request] L7: %d", baseCtx.L7)

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionNext()
	}

	extractor := Extractor{Payload: payload}
	attrs, trace, err := extractor.do()
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != GRPC || baseCtx.DstPort != 9000 {
		return sdk.ActionNext()
	}
	sdk.Warn("[Response] L7: %d", baseCtx.L7)

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}

	extractor := Extractor{Payload: payload}
	attrs, trace, err := extractor.do()
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpRespActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnCheckPayload(baseCtx *sdk.ParseCtx) (uint8, string) {
	return 0, ""
}

func (p httpHook) OnParsePayload(baseCtx *sdk.ParseCtx) sdk.Action {
	return sdk.ActionNext()
}

func main() {
	sdk.Warn("wasm register http hook")
	sdk.SetParser(httpHook{})
}
