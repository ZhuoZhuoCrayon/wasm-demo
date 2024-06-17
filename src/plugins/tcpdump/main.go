package main

import (
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common"
	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	// 将 nottinygc 作为 TinyGo 编译 WASI 的一个替代内存分配器，默认的内存分配器在数据量大的场景会有性能问题
	_ "github.com/wasilibs/nottinygc"
)

var dstPorts = []uint16{17216, 11048, 16830}

type httpHook struct {
}

func (p httpHook) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_PAYLOAD_PARSE,
	}
}

func (p httpHook) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	return sdk.ActionNext()
}

func (p httpHook) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	return sdk.ActionNext()
}

func (p httpHook) OnCheckPayload(baseCtx *sdk.ParseCtx) (uint8, string) {
	if !common.IsPortMatched(baseCtx.DstPort, dstPorts) {
		return 0, ""
	}
	sdk.Info("[OnCheckPayload] %v:%v -> %v:%v", baseCtx.SrcIP, baseCtx.SrcPort, baseCtx.DstIP, baseCtx.DstPort)
	payload, _ := baseCtx.GetPayload()
	sdk.Info("[OnCheckPayload] payload -> %x", payload)
	return 0, ""
}

func (p httpHook) OnParsePayload(baseCtx *sdk.ParseCtx) sdk.Action {
	return sdk.ActionNext()
}

func main() {
	sdk.Warn("wasm register plugin -> tcpdump")
	sdk.SetParser(httpHook{})
}
