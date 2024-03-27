# Wasm demo



## 项目结构


```text
.
├── go.mod
└── src
    ├── plugins                        // 存放 WASM 插件
    │   ├── grpc_data_extractor        // gRPC body 提取插件
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
    │   ├── http2_header_extractor   // HTTP2 Header 提取插件
    │   │   ├── Makefile
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── http2_header_extractor.wasm
    │   │   ├── main.go
    │   │   └── readme.md
    │   └── utils                  //  公共代码
    │       └── http2              //  HTTP2 Payload 解析模块
    │           ├── go.mod
    │           ├── http2.go
    │           └── http2_test.go
    └── userinfoservice           // gRPC 服务 demo
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
