package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common"
	uhttp "github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"io"
	"net/http"
	"strconv"
)

const (
	HTTP1        uint8 = 20
	MinBodyBytes       = 19
)

var (
	ExtractError = errors.New("failed to extract data")
)

var expectDataFields = map[string]bool{
	"code":      true,
	"result":    true,
	"message":   true,
	"code_name": true,
}

var expectFieldPatterns = map[string]string{
	"code":      `"code":\s*"?(\d+)"?`,
	"result":    `"result":\s*([a-zA-Z]+)`,
	"message":   `"message":\s*"([^"]+)"`,
	"code_name": `"code_name":\s*"([^"]+)"`,
}

func getOrDefault(kv map[string]string, field string, defaultValue string) string {
	var ok bool
	var val string
	if val, ok = kv[field]; !ok {
		return defaultValue
	}
	return val
}

func formatCode(codeStr string) int32 {
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		code = 500
	}
	// 0 转为 200
	if code == 0 {
		code = 200
	}
	return int32(code)
}

func getStatus(code int32) sdk.RespStatus {
	if code != 0 && code != 200 {
		return sdk.RespStatusServerErr
	}
	return sdk.RespStatusOk
}

func formatKv(kv map[string]string) ([]sdk.KeyVal, *sdk.Response, error) {

	if kv == nil || len(kv) == 0 {
		return nil, nil, ExtractError
	}
	sdk.Warn("[formatKv] extract kv -> %v", kv)
	attrs := make([]sdk.KeyVal, 0, len(kv))
	for k, v := range kv {
		attrs = append(attrs, sdk.KeyVal{Key: k, Val: v})
	}

	// code：0 / 200 表示正常
	// result：可能存在
	// message：可能存在
	// code_name：可能存在

	codeStr, ok := kv["code"]
	if !ok {
		return nil, nil, ExtractError
	}

	code := formatCode(codeStr)
	status := getStatus(code)

	// handle Error
	if status == sdk.RespStatusServerErr {
		message := getOrDefault(kv, "message", "")
		codeName := getOrDefault(kv, "code_name", "ERROR")

		sdkResp := &sdk.Response{
			Code:      &code,
			Status:    &status,
			Result:    message,
			Exception: fmt.Sprintf("%v(%v)", codeName, message),
		}

		return attrs, sdkResp, nil
	}

	sdkResp := &sdk.Response{Code: &code, Status: &status}
	return attrs, sdkResp, nil
}

type httpHook struct {
	sdk.DefaultParser
}

func (p httpHook) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		// sdk.HOOK_POINT_HTTP_REQ,
		sdk.HOOK_POINT_HTTP_RESP,
	}
}

func (p httpHook) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	return sdk.ActionNext()
}

func extraBodyFromResp(resp *http.Response) (body []byte, err error) {
	defer resp.Body.Close()

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		g, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}

		defer g.Close()

		body, err = io.ReadAll(g)
		if err != nil {
			return nil, err
		}

	default:
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	baseCtx := &ctx.BaseCtx
	if baseCtx.L7 != HTTP1 {
		return sdk.ActionNext()
	}

	if baseCtx.BufSize < MinBodyBytes {
		return sdk.ActionNext()
	}

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionNext()
	}
	// sdk.Warn("[payload] %v", string(payload))

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(payload)), &http.Request{})
	if err != nil {
		return sdk.ActionNext()
	}
	// sdk.Warn("[header] %v", resp.Header)

	xBkapiRequestId := resp.Header.Get("X-Bkapi-Request-ID")
	if len(xBkapiRequestId) == 0 {
		return sdk.ActionNext()
	}

	contentType := resp.Header.Get("Content-Type")
	if !common.IsPrefixMatched(contentType, []string{uhttp.JSON, uhttp.PlainText}) {
		return sdk.ActionNext()
	}

	body, err := extraBodyFromResp(resp)
	if err != nil {
		return sdk.ActionNext()
	}
	if len(body) < MinBodyBytes {
		return sdk.ActionNext()
	}
	// sdk.Warn("[body] %v", string(body))

	kv, err := uhttp.ExtractBytes(body, uhttp.JSON, expectDataFields)
	if err != nil || len(kv) == 0 {
		kv, err = uhttp.ExtractIncompleteJSON(body, expectFieldPatterns)
		if err != nil {
			return sdk.ActionNext()
		}
	}

	attrs, sdkResp, err := formatKv(kv)
	if err != nil {
		return sdk.ActionNext()
	}

	return sdk.HttpRespActionAbortWithResult(sdkResp, &sdk.Trace{}, attrs)
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
