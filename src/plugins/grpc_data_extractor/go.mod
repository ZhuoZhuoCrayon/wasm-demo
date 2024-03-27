module github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/grpc_data_extractor

go 1.20

require (
	github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2 v0.0.0-00010101000000-000000000000
	github.com/deepflowio/deepflow-wasm-go-sdk v0.0.0-20240228064431-d05daaa99d08
	google.golang.org/protobuf v1.33.0
)

replace github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/utils/http2 => ../utils/http2
