module github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/tcpdump

go 1.20

require (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common v0.0.0-00010101000000-000000000000
	github.com/deepflowio/deepflow-wasm-go-sdk v0.0.0-20240228064431-d05daaa99d08
	github.com/wasilibs/nottinygc v0.7.1
)

require github.com/magefile/mage v1.14.0 // indirect

replace github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/common => ../utils/common
