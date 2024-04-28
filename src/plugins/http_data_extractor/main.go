package main

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common"
	uhttp "github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	// 将 nottinygc 作为 TinyGo 编译 WASI 的一个替代内存分配器，默认的内存分配器在数据量大的场景会有性能问题
	_ "github.com/wasilibs/nottinygc"
	"net/http"
)

const (
	HTTP1 uint8 = 20
)

var (
	ExtractError = errors.New("failed to extract data")
)

var httpPathPrefixes = []string{
	"/v3/r/mpay",
	"/api/checkout",
}

var expectDataFields = map[string]bool{
	"openid": true,
	// "err_code": true,
	// "ret":      true,
	"userId": true,
	// "orderId": true,
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
	trace := &sdk.Trace{}
	return attrs, trace, nil
}

type httpHook struct {
}

func (p httpHook) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
		// sdk.HOOK_POINT_HTTP_RESP,
	}
}

func (p httpHook) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != HTTP1 || !common.IsPrefixMatched(ctx.Path, httpPathPrefixes) {
		return sdk.ActionNext()
	}
	// sdk.Info("[OnHttpReq] path -> %v", ctx.Path)
	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	// sdk.Info("[OnHttpReq] payload -> %#v", payload)
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		return sdk.ActionNext()
	}
	contentType := req.Header.Get("Content-Type")
	// sdk.Info("[OnHttpReq] header -> %v, Content-Type -> %v", req.Header, contentType)
	kv, err := uhttp.Extract(req.Body, contentType, expectDataFields)
	if err != nil {
		return sdk.ActionNext()
	}
	attrs, trace, err := formatKv(kv)
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	// TODO(crayon) Response 计算量太大，暂不开启
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
