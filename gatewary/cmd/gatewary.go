package main

import (
	"log"

	"github.com/dollarkillerx/smart-dns-go/gatewary/core"
	"github.com/dollarkillerx/smart-dns-go/gatewary/define"
)

func main() {
	if define.BaseConfig.Debug {
		log.SetFlags(log.Llongfile | log.LstdFlags)
	}

	app := core.Core{}
	log.Println("Gatewary Run: ", define.BaseConfig.ListenAddr)
	if err := app.Core(); err != nil {
		log.Fatalln(err)
	}
}
