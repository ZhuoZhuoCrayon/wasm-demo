plugin_name ?= http_resp_extractor
compile:
	tinygo build -o ${plugin_name}.wasm -gc=precise -target=wasi -panic=trap -scheduler=none -no-debug ./main.go

apply:
	deepflow-ctl plugin delete ${plugin_name}
	deepflow-ctl plugin create --type wasm --image ${plugin_name}.wasm --name ${plugin_name}
	deepflow-ctl plugin list | grep ${plugin_name}
