module github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/http2_header_extractor

go 1.20

require (
	github.com/deepflowio/deepflow-wasm-go-sdk v0.0.0-20240228064431-d05daaa99d08
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2 v0.0.0-00010101000000-000000000000
	github.com/wasilibs/nottinygc v0.7.1
	golang.org/x/net v0.9.0
)

require github.com/magefile/mage v1.14.0 // indirect

replace (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2 => ../utils/http2
)
