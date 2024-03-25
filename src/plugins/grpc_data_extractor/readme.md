# UserInfoService

## QuickStart

### Init

```shell
go mod tidy
```


generate pb 

```shell
midir protos

# protoc --go_out=protos --go_opt=paths=source_relative userinfo.proto
# refer: 
protoc --go-plugin_out=./protos --go-plugin_opt=paths=source_relative ./userinfo.proto
```

### Compile

```shell
# 使用 nottinygc 替换 TinyGo 原来的内存分配器需要增加编译参数：-gc=custom 和 -tags=custommalloc
tinygo build -o grpc_data_extractor.wasm -target=wasi -panic=trap -scheduler=none -no-debug ./main.go
tinygo build -o grpc_data_extractor.wasm -target wasi -gc=precise -panic=trap -scheduler=none -no-debug *.go
```

#### Apply

```shell
deepflow-ctl plugin create --type wasm --image grpc_data_extractor.wasm --name grpc_data_extractor
```
