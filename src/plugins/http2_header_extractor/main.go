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
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	// 将 nottinygc 作为 TinyGo 编译 WASI 的一个替代内存分配器，默认的内存分配器在数据量大的场景会有性能问题
	_ "github.com/wasilibs/nottinygc"
)

const (
	GRPC uint8 = 41
)

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

	sdk.Info("L7: %d", baseCtx.L7)

	if baseCtx.L7 != GRPC {
		return sdk.ActionNext()
	}

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	sdk.Warn("payload(ascii): %s", payload)
	sdk.Warn("payload: %x", payload)

	attr := []sdk.KeyVal{
		{
			Key: "openid",
			Val: "123456789",
		},
	}

	trace := &sdk.Trace{
		TraceID: "123456789",
		SpanID:  "123456788",
	}

	return sdk.HttpReqActionAbortWithResult(nil, trace, attr)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	return sdk.ActionNext()
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
