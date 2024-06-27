package main

import (
	"fmt"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/propertyf"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/statf"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/collector/tars/property"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/collector/tars/stat"
	"os"
)

type Adapter struct {
	Name     string
	Servant  string
	Endpoint string
}

type ServerConfig struct {
	App      string
	Server   string
	LogPath  string
	LogLevel string
	Adapters []Adapter
}

func InitConfig() ServerConfig {
	statObjAdapter := Adapter{
		Name:     "collector.tars.StatObjAdapter",
		Servant:  "collector.tars.StatObj",
		Endpoint: "tcp -h 127.0.0.1 -p 10892 -t 60000",
	}
	propertyObjAdapter := Adapter{
		Name:     "collector.tars.PropertyObjAdapter",
		Servant:  "collector.tars.PropertyObj",
		Endpoint: "tcp -h 127.0.0.1 -p 10891 -t 60000",
	}
	cfg := ServerConfig{
		App:      "collector",
		Server:   "tars",
		LogPath:  "/tmp/log/gse2_bkte/bk-collector/tars/",
		LogLevel: "DEBUG",
		Adapters: []Adapter{statObjAdapter, propertyObjAdapter},
	}
	return cfg
}

func GenConfig(configPath string) error {
	cfg := InitConfig()

	xmlString := fmt.Sprintf(`<tars>
    <application>
        <server>
            app=%s
            server=%s
            logpath=%s
            logLevel=%s
`, cfg.App, cfg.Server, cfg.LogPath, cfg.LogLevel)

	for _, adapter := range cfg.Adapters {
		xmlString += fmt.Sprintf(`            <%s>
                endpoint=%s
                handlegroup=%s
                servant=%s
            </%s>
`, adapter.Name, adapter.Endpoint, adapter.Name, adapter.Servant, adapter.Name)
	}

	xmlString += `        </server>
    </application>
</tars>`

	err := os.WriteFile(configPath, []byte(xmlString), 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Get server config
	configPath := "config/config.auto.conf"
	if err := GenConfig(configPath); err != nil {
		fmt.Printf("config init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	tars.ServerConfigPath = configPath
	cfg := tars.GetServerConfig()

	statApp := new(statf.StatF)
	// New servant imp
	statImp := stat.NewStatImp()
	err := statImp.Init()
	if err != nil {
		fmt.Printf("statImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	tars.AddServantWithContext(statApp, statImp, cfg.App+"."+cfg.Server+".StatObj")

	propertyApp := new(propertyf.PropertyF)
	propertyImp := property.NewPropertyImp()
	err = propertyImp.Init()
	if err != nil {
		fmt.Printf("propertyImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	tars.AddServantWithContext(propertyApp, propertyImp, cfg.App+"."+cfg.Server+".PropertyObj")

	// Run application
	tars.Run()
}
