package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"

	// 将 nottinygc 作为 TinyGo 编译 WASI 的一个替代内存分配器，默认的内存分配器在数据量大的场景会有性能问题
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"github.com/valyala/fastjson"
	_ "github.com/wasilibs/nottinygc"
	"net/http"
	"strings"
)

const (
	HTTP1 uint8 = 20
)

var httpPathPrefixes = []string{
	"/v3/r/mpay",
}

var expectDataFields = map[string]bool{
	"openid":   true,
	"err_code": true,
	"ret":      true,
}

func isPathMatched(path string, prefixes []string) bool {
	sdk.Info("[isPathMatched] path -> %s", path)
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			sdk.Info("[isPathMatched] path -> %s, matched -> %s", path, prefix)
			return true
		}
	}
	return false
}

func extract(b io.ReadCloser, fields map[string]bool) ([]sdk.KeyVal, *sdk.Trace, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		return nil, nil, err
	}
	sdk.Info("[extract] body -> %#v", body)
	fv, err := fastjson.ParseBytes(body)
	if err != nil {
		return nil, nil, err
	}
	attrs := make([]sdk.KeyVal, 0, len(expectDataFields))
	for field := range fields {
		if !fv.Exists(field) {
			continue
		}
		v := ""
		switch fv.Get(field).Type() {
		case fastjson.TypeString:
			v = string(fv.GetStringBytes(field))
		case fastjson.TypeNumber:
			v = strconv.Itoa(fv.GetInt(field))
		}
		if v != "" {
			attrs = append(attrs, sdk.KeyVal{Key: field, Val: string(v)})
		}
	}
	sdk.Info("[extract] attrs -> %v", attrs)
	return attrs, nil, nil
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
	if ctx.BaseCtx.L7 != HTTP1 || !isPathMatched(ctx.Path, httpPathPrefixes) {
		return sdk.ActionNext()
	}
	payload, err := ctx.BaseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		return sdk.ActionNext()
	}
	attrs, trace, err := extract(req.Body, expectDataFields)
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	if ctx.BaseCtx.L7 != HTTP1 {
		return sdk.ActionNext()
	}
	payload, err := ctx.BaseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(payload)), nil)
	if err != nil {
		return sdk.ActionNext()
	}
	attrs, trace, err := extract(resp.Body, expectDataFields)
	if err != nil {
		return sdk.ActionNext()
	}
	return sdk.HttpReqActionAbortWithResult(nil, trace, attrs)
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
