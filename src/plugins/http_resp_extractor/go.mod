module github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/http_resp_extractor

go 1.20

require (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http v0.0.0-00010101000000-000000000000
	github.com/deepflowio/deepflow-wasm-go-sdk v0.0.0-20240511132801-1f9ec0e0706e
)

require (
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common => ../utils/common
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http => ../utils/http
)
