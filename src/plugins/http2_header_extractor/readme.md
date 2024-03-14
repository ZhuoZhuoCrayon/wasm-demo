# UserInfoService

## QuickStart

### Init

```shell
go mod tidy
```

### Compile

```shell
# 使用 nottinygc 替换 TinyGo 原来的内存分配器需要增加编译参数：-gc=custom 和 -tags=custommalloc
tinygo build -o http2_header_extractor.wasm -gc=custom -tags=custommalloc -target=wasi -panic=trap -scheduler=none -no-debug ./main.go
```

#### Apply

```shell
deepflow-ctl plugin create  --type wasm --image http2_header_extractor.wasm --name http2_header_extractor
```
