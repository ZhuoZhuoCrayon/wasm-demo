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
	"errors"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2"
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
	GRPC uint8  = 41
	Port uint16 = 9000
)

var (
	ExtractError = errors.New("failed to extract headers")
)

var ExpectedHeaderFields = map[string]bool{
	"openid": true,
}

func formatKv(kv map[string]string) ([]sdk.KeyVal, *sdk.Trace, error) {
	if kv == nil || len(kv) == 0 {
		return nil, nil, ExtractError
	}
	sdk.Info("[formatKv] extract kv -> %v", kv)
	attrs := make([]sdk.KeyVal, 0, len(kv))
	for k, v := range kv {
		attrs = append(attrs, sdk.KeyVal{Key: k, Val: v})
	}
	trace := &sdk.Trace{
		TraceID: attrs[0].Val,
	}
	return attrs, trace, nil
}

type http2HookCtx struct {
	isResp       bool
	kv           map[string]string
	hPackDecoder *hpack.Decoder
}

func (h *http2HookCtx) onHeader(fh http2.FrameHeader, bf []byte) (err error) {
	sdk.Info("[onHeader] BlockFragment -> %x, IsResp -> %v", bf[5:], h.isResp)
	if _, err = h.hPackDecoder.Write(bf); err != nil {
		return err
	}
	return nil
}

func newHttp2HookCtx(isResp bool) *http2HookCtx {
	hc := &http2HookCtx{isResp: isResp, kv: make(map[string]string, len(ExpectedHeaderFields)), hPackDecoder: nil}
	hc.hPackDecoder = hpack.NewDecoder(4096, func(hf hpack.HeaderField) {
		if _, ok := ExpectedHeaderFields[hf.Name]; ok {
			hc.kv[hf.Name] = hf.Value
		}
	})
	return hc
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
	if baseCtx.L7 != GRPC || baseCtx.DstPort != Port {
		return sdk.ActionNext()
	}
	sdk.Warn("[OnHttpReq] L7: %d", baseCtx.L7)
	payload, err := baseCtx.GetPayload()
	sdk.Warn("[OnHttpReq] Payload -> %v", payload)
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	hc := newHttp2HookCtx(false)
	parser := http2.NewParser()
	parser.RegisterHook(http2.FrameHeaders, hc.onHeader)
	parser.RegisterHook(http2.FrameContinuation, hc.onHeader)
	parser.Do(payload)
	attrs, trace, err := formatKv(hc.kv)
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != GRPC || baseCtx.SrcPort != Port {
		return sdk.ActionNext()
	}
	sdk.Warn("[OnHttpResp] L7: %d", baseCtx.L7)
	payload, err := baseCtx.GetPayload()
	sdk.Warn("[OnHttpResp] Payload -> %v", payload)
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	hc := newHttp2HookCtx(true)
	parser := http2.NewParser()
	parser.RegisterHook(http2.FrameHeaders, hc.onHeader)
	parser.RegisterHook(http2.FrameContinuation, hc.onHeader)
	parser.Do(payload)
	attrs, trace, err := formatKv(hc.kv)
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
