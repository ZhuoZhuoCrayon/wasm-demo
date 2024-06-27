package main

import (
	"context"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/TestApp/HelloGo/TestApp"
)

const (
	// NodeIp       string = "hellogo-service"
	NodeIp       string = "127.0.0.1"
	RegistryPort int    = 30890
	// ServerPort   int    = 30001
	ServerPort int = 10015
)

func main() {
	comm := tars.NewCommunicator()
	// Server 已注册到 registry
	//obj := fmt.Sprintf("TestApp.HelloGo.SayHelloObj")
	// comm.SetProperty("locator", fmt.Sprintf("tars.tarsregistry.QueryObj@tcp -h %s -p %d", NodeIp, RegistryPort))
	obj := fmt.Sprintf("TestApp.HelloGo.SayHelloObj@tcp -h %s -p %d -t 60000", NodeIp, ServerPort)
	app := new(TestApp.SayHello)
	comm.StringToProxy(obj, app)
	var out, i int32
	i = 128888

	ctx := context.Background()
	c := make(map[string]string)
	c["a"] = "b"
	ret, err := app.AddWithContext(ctx, i, i*2, &out, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
