package main

import (
	"context"
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/current"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/TestApp/HelloGo/TestApp"
	"os"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()
	// New properties
	p := NewProperties()
	// New servant imp
	imp := NewSayHelloImp(p)
	err := imp.Init()
	if err != nil {
		fmt.Printf("SayHelloImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(TestApp.SayHello)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".SayHelloObj")

	ctx := context.Background()
	tokens := map[string]string{"a": "b", "b": "c"}
	current.SetRequestContext(ctx, tokens)
	// Run application
	tars.Run()
}
