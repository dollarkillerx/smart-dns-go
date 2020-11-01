package main

import (
	"log"

	"github.com/dollarkillerx/smart-dns-go/gatewary/core"
	"github.com/dollarkillerx/smart-dns-go/gatewary/pkg/config"
)

func main() {
	if config.BaseConfig.Debug {
		log.SetFlags(log.Llongfile | log.LstdFlags)
	}

	app := core.NewCore()
	log.Println("Gatewary Run: ", config.BaseConfig.ListenAddr)
	if err := app.Core(); err != nil {
		log.Fatalln(err)
	}
}
