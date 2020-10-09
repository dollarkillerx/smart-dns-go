package main

import (
	"log"
	"os"

	"github.com/dollarkillerx/smart-dns-go/gatewary/core"
	"github.com/dollarkillerx/smart-dns-go/gatewary/pkg/config"
)

func main() {
	os.Setenv("GODEBUG", "x509ignoreCN=0")

	if config.BaseConfig.Debug {
		log.SetFlags(log.Llongfile | log.LstdFlags)
	}

	app := core.NewCore()
	log.Println("Gatewary Run: ", config.BaseConfig.ListenAddr)
	if err := app.Core(); err != nil {
		log.Fatalln(err)
	}
}
