module github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/http_data_extractor

go 1.20

require (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common v0.0.0-00010101000000-000000000000
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http v0.0.0-00010101000000-000000000000
	github.com/deepflowio/deepflow-wasm-go-sdk v0.0.0-20240417023450-0d59087dd02f
	github.com/wasilibs/nottinygc v0.7.1
)

require (
	github.com/magefile/mage v1.14.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common => ../utils/common
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http => ../utils/http
)
