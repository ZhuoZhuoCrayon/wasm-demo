package main

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/TestApp/HelloGo/TestApp"
)

const (
	NodeIp       string = "localhost"
	RegistryPort int    = 30890
	ServerPort   int    = 30001
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
	i = 124
	ret, err := app.Add(i, i*2, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
