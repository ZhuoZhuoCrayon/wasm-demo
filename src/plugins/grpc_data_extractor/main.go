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
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/grpc_data_extractor/protos"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"reflect"
	"strconv"
)

type FrameType uint8

type Flags uint8

func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

const (
	GRPC           uint8     = 41
	FrameHeaderLen int       = 9
	FrameTypeMax   FrameType = 0x9
	FrameData      FrameType = 0x0
	FlagDataPadded Flags     = 0x8
	FlagEndStream  Flags     = 0x1
)

var (
	ParsePayloadError = errors.New("failed to l7 payload")
	ExtractError      = errors.New("failed to extract data")
)

var ExpectedDataFields = map[string]bool{
	"Openid": true,
	"UserId": true,
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

func readDataFramePayload(p []byte, fh FrameHeader) (remain []byte, dp []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	var padLength uint8
	if fh.Flags.Has(FlagDataPadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, nil, err
		}
	} else {
		padLength = 0
	}
	if len(p)-int(padLength) < 0 || int(padLength) > int(fh.Length) {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:int(fh.Length)-int(padLength)], nil
}

func readFramePayload(p []byte, fh FrameHeader) (remain []byte, fp []byte, err error) {
	if isOverflow := frameOverflow(p, fh); isOverflow {
		return nil, nil, ParsePayloadError
	}
	return p[fh.Length:], p[:fh.Length], nil
}

func checkAndAddFieldsToMap(s interface{}, fields map[string]bool, kv map[string]string) {
	v := reflect.ValueOf(s).Elem()
	for field := range fields {
		f := v.FieldByName(field)
		if f.IsValid() {
			switch f.Kind() {
			case reflect.String:
				kv[field] = f.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				kv[field] = strconv.FormatInt(f.Int(), 10)
			}
		}
	}
}

// parse 解析 HTTP2 payload
func parse(payload []byte, isResp bool) (kv map[string]string, err error) {
	sdk.Warn("start to l7 Payload: %x", payload)

	p := payload
	var fh FrameHeader
	kv = make(map[string]string, len(ExpectedDataFields))
	// 输入 wasm 插件的 Payload 可能存在数据截断的情况，使用 break 结束循环而不是 return nil, err 来尽可能匹配前面符合规则的帧
	for len(p) > FrameHeaderLen {
		if p, fh, err = readFrameHeader(p); err != nil {
			break
		}
		sdk.Info("[l7] FrameType -> %#x", fh.Type)
		if fh.Type == FrameData {
			var dataPayload []byte
			if p, dataPayload, err = readDataFramePayload(p, fh); err != nil {
				break
			}
			sdk.Info("[l7] DataPayload -> %x, isResp -> %v", dataPayload[5:], isResp)
			if isResp {
				resp := &pb.UserInfo{}
				if err = resp.UnmarshalVT(dataPayload[5:]); err != nil {
					sdk.Warn("[l7] failed to unmarshal resp payload")
					break
				}
				sdk.Info("[l7] resp -> %v", resp)
				checkAndAddFieldsToMap(resp, ExpectedDataFields, kv)
			} else {
				req := &pb.GetUserInfoRequest{}
				if err = req.UnmarshalVT(dataPayload[5:]); err != nil {
					sdk.Warn("[l7] failed to unmarshal req payload")
					break
				}
				sdk.Info("[l7] req -> %v", req)
				checkAndAddFieldsToMap(req, ExpectedDataFields, kv)
			}
		} else {
			if p, _, err = readFramePayload(p, fh); err != nil {
				break
			}
		}
	}
	return kv, nil
}

type Extractor struct {
	IsResp  bool
	Payload []byte
}

func (e *Extractor) do() ([]sdk.KeyVal, *sdk.Trace, error) {
	kv, err := parse(e.Payload, e.IsResp)
	if err != nil || kv == nil {
		return nil, nil, ExtractError
	}
	sdk.Info("extract kv: %v", kv)
	if len(kv) == 0 {
		return nil, nil, ExtractError
	}

	attrs := make([]sdk.KeyVal, 0, len(kv))
	for k, v := range kv {
		attrs = append(attrs, sdk.KeyVal{Key: k, Val: v})
	}

	trace := &sdk.Trace{}
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

	extractor := Extractor{Payload: payload, IsResp: false}
	attrs, trace, err := extractor.do()
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != GRPC || baseCtx.SrcPort != 9000 {
		return sdk.ActionNext()
	}
	sdk.Warn("[Response] L7: %d", baseCtx.L7)

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}

	extractor := Extractor{Payload: payload, IsResp: true}
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
