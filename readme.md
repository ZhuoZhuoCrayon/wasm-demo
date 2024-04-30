# Wasm demo



## 项目结构


```text
├── go.mod
├── readme.md
└── src
├── core
│   ├── ftrpc                    // tRPC 解析 demo
│   │   ├── Makefile
│   │   ├── frpc_test.go
│   │   ├── ftrpc.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── side
│   │   │   └── main.go
│   │   └── trpc
│   │       ├── trpc.pb.go
│   │       └── trpc.proto
│   └── query                   // rpc demo 依赖模块
│       ├── go.mod
│       ├── query.go
│       └── query_test.go
├── plugins
│   ├── grpc_data_extractor    // gRPC DATA 提取
│   │   ├── Makefile
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── grpc_data_extractor.wasm
│   │   ├── main.go
│   │   ├── protos
│   │   │   ├── userinfo.pb.go
│   │   │   └── userinfo_vtproto.pb.go
│   │   ├── readme.md
│   │   ├── tests
│   │   │   └── protos_test.go
│   │   └── userinfo.proto
│   ├── http2_header_extractor      // HTTP/2 Headers 提取
│   │   ├── Makefile
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── http2_header_extractor.wasm
│   │   ├── main.go
│   │   └── readme.md
│   ├── http_data_extractor         // HTTP DATA 提取
│   │   ├── Makefile
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── http_data_extractor.wasm
│   │   └── main.go
│   └── utils
│       ├── common
│       │   ├── common.go
│       │   ├── common_test.go
│       │   ├── go.mod
│       │   └── go.sum
│       ├── http                // HTTP 协议解析
│       │   ├── go.mod
│       │   ├── go.sum
│       │   ├── http.go
│       │   └── http_test.go
│       └── http2               // HTTP/2 协议解析
│           ├── go.mod
│           ├── http2.go
│           └── http2_test.go
├── timeseriesqueryservice      // gRPC Stream demo
│   ├── Makefile
│   ├── client
│   │   ├── bstream
│   │   │   └── client.go
│   │   ├── client.go
│   │   ├── cstream
│   │   │   └── client.go
│   │   └── sstream
│   │       └── client.go
│   ├── go.mod
│   ├── go.sum
│   ├── server
│   │   └── server.go
│   └── timeseriesquery
│       ├── time_series_query.pb.go
│       ├── time_series_query.proto
│       └── time_series_query_grpc.pb.go
├── trpc_timeseriesqueryservice             // tRPC demo
│   ├── Makefile
│   ├── client
│   │   ├── bstream
│   │   │   └── client.go
│   │   ├── client.go
│   │   ├── cstream
│   │   │   └── client.go
│   │   ├── sstream
│   │   │   └── client.go
│   │   ├── trpc_go.yaml
│   │   └── ustream
│   │       └── client.go
│   ├── conf
│   │   └── trpc_go.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── greeter.tar.gz
│   ├── main.go
│   ├── stke_build
│   │   ├── README.md
│   │   ├── build.sh
│   │   ├── env
│   │   ├── haha.yaml
│   │   ├── image
│   │   │   └── Dockerfile
│   │   ├── init-env.sh
│   │   └── script
│   │       └── entrypoint.sh
│   ├── target
│   │   └── timeseriesquery
│   │       └── timeseriesquery
│   ├── timeseriesquery
│   │   ├── time_series_query.pb.go
│   │   ├── time_series_query.proto
│   │   └── time_series_query.trpc.go
│   └── timeseriesquery.tar.gz
└── userinfoservice                     // gRPC Unary demo
├── client
│   └── main.go
├── go.mod
├── go.sum
├── main.go
├── protos
│   └── userinfo
│       ├── userinfo.pb.go
│       └── userinfo_grpc.pb.go
├── readme.md
└── userinfo.proto
```