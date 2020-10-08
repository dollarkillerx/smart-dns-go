package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	ListenAddr string `json:"listen_addr" yaml:"ListenAddr"`
	CoreAddr   string `json:"core_addr" yaml:"CoreAddr"`
	Debug      bool   `json:"debug" yaml:"Debug"`
	PublicDNS  string `json:"public_dns" yaml:"PublicDNS"`

	// Auth
	SSLPem  string `json:"pem" yaml:"SSLPem"`
	SSLAuth string `json:"ssl_auth" yaml:"SSLAuth"`
}

var BaseConfig = initBase()

func initBase() *config {
	file, err := ioutil.ReadFile("./gatewary/configs/config.yaml")
	if err != nil {
		if err := ioutil.WriteFile("./gatewary/configs/config.yaml", []byte(cfp), 00666); err != nil {
			log.Fatalln(err)
		} else {
			log.Fatalln("Please fill out the profile")
		}
	}
	var cfg config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

const cfp = `
ListenAddr: "0.0.0.0:53"
CoreAddr: "0.0.0.0:8081"
Debug: true
PublicDNS: "8.8.8.8:53"

SSLPem: "./configs/s1.key"
SSLAuth: "key_password"
`
