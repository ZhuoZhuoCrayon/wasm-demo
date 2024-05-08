package main

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/TestApp/HelloGo/TestApp"
	"os"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(SayHelloImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("SayHelloImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(TestApp.SayHello)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".SayHelloObj")

	// Run application
	tars.Run()
}
