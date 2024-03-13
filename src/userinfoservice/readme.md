# UserInfoService

## QuickStart

### Init

```shell
go mod tidy
```

### Generate Code

```shell
protoc --go_out=protos/userinfo --go_opt=paths=source_relative --go-g
rpc_out=protos/userinfo --go-grpc_opt=paths=source_relative userinfo.proto
```

#### Run Server

```shell
go run main.go
```

#### Run Client

```shell
go run client/main.go
```
