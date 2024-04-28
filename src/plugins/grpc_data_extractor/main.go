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
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/grpc_data_extractor/protos"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"reflect"
	"strconv"
)

const (
	GRPC uint8  = 41
	Port uint16 = 9000
)

var (
	ExtractError = errors.New("failed to extract data")
)

var ExpectedDataFields = map[string]bool{
	"Openid": true,
	"UserId": true,
}

var FieldMap = map[string]string{
	"Openid": "openid",
	"UserId": "userId",
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

func formatKv(kv map[string]string) ([]sdk.KeyVal, *sdk.Trace, error) {
	if kv == nil || len(kv) == 0 {
		return nil, nil, ExtractError
	}
	sdk.Info("[formatKv] extract kv -> %v", kv)
	attrs := make([]sdk.KeyVal, 0, len(kv))
	for k, v := range kv {
		field, isOk := FieldMap[k]
		if !isOk {
			field = k
		}
		attrs = append(attrs, sdk.KeyVal{Key: field, Val: v})
	}
	trace := &sdk.Trace{}
	return attrs, trace, nil
}

type http2HookCtx struct {
	isResp bool
	kv     map[string]string
}

func (h *http2HookCtx) onData(fh http2.FrameHeader, bf []byte) (err error) {
	if len(bf) < 5 {
		sdk.Warn("[onData] Invalid gRPC DATA")
		return ExtractError
	}
	sdk.Info("[onData] BlockFragment -> %x, IsResp -> %v", bf[5:], h.isResp)
	if h.isResp {
		resp := &pb.UserInfo{}
		if err = resp.UnmarshalVT(bf[5:]); err != nil {
			sdk.Warn("[onData] failed to unmarshal resp payload")
			return err
		}
		sdk.Info("[onData] resp -> %v", resp)
		checkAndAddFieldsToMap(resp, ExpectedDataFields, h.kv)
	} else {
		req := &pb.GetUserInfoRequest{}
		if err = req.UnmarshalVT(bf[5:]); err != nil {
			sdk.Warn("[onData] failed to unmarshal req payload")
			return err
		}
		sdk.Info("[onData] req -> %v", req)
		checkAndAddFieldsToMap(req, ExpectedDataFields, h.kv)
	}
	return nil
}

type httpHook struct {
	sdk.DefaultParser
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
	sdk.Warn("[OnHttpReq] L7 -> %d", baseCtx.L7)
	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionNext()
	}
	sdk.Warn("[OnHttpReq] Payload -> %v", payload)
	hc := &http2HookCtx{isResp: false, kv: make(map[string]string, len(ExpectedDataFields))}
	parser := http2.NewParser()
	parser.RegisterHook(http2.FrameData, hc.onData)
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
	sdk.Warn("[OnHttpResp] L7 -> %d", baseCtx.L7)
	payload, err := baseCtx.GetPayload()
	sdk.Warn("[OnHttpResp] Payload -> %v", payload)
	if err != nil {
		return sdk.ActionNext()
	}
	hc := &http2HookCtx{isResp: true, kv: make(map[string]string, len(ExpectedDataFields))}
	parser := http2.NewParser()
	parser.RegisterHook(http2.FrameData, hc.onData)
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
