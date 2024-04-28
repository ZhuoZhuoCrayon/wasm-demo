plugin_name ?= grpc_data_extractor
pb:
	# refer: https://github.com/knqyf263/go-plugin
	# protobuf/proto 在 WASM 插件内执行会出现报错，故使用 go-plugin 生成 pb 文件
	mkdir -p ./protos && protoc --go-plugin_out=./protos --go-plugin_opt=paths=source_relative ./userinfo.proto

compile:
	tinygo build -o ${plugin_name}.wasm -target wasi -gc=precise -panic=trap -scheduler=none -no-debug *.go

apply:
	deepflow-ctl plugin delete ${plugin_name}
	deepflow-ctl plugin create --type wasm --image ${plugin_name}.wasm --name ${plugin_name}
	deepflow-ctl plugin list | grep ${plugin_name}
