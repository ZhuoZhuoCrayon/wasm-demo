compile:
	# 使用 nottinygc 替换 TinyGo 原来的内存分配器需要增加编译参数：-gc=custom 和 -tags=custommalloc
	tinygo build -o http_data_extractor.wasm -gc=custom -tags=custommalloc -target=wasi -panic=trap -scheduler=none -no-debug ./main.go

apply:
	deepflow-ctl plugin create --type wasm --image http_data_extractor.wasm --name http_data_extractor
