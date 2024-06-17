plugin_name ?= tcpdump
compile:
	# 使用 nottinygc 替换 TinyGo 原来的内存分配器需要增加编译参数：-gc=custom 和 -tags=custommalloc
	tinygo build -o ${plugin_name}.wasm -gc=custom -tags=custommalloc -target=wasi -panic=trap -scheduler=none -no-debug ./main.go

apply:
	deepflow-ctl plugin delete ${plugin_name}
	deepflow-ctl plugin create --type wasm --image ${plugin_name}.wasm --name ${plugin_name}
	deepflow-ctl plugin list | grep ${plugin_name}
