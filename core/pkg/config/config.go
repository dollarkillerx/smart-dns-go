package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

type config struct {
	ListenAddr string `json:"listen_addr" yaml:"ListenAddr"`
	Debug      bool   `json:"debug" yaml:"Debug"`
	PublicDNS  string `json:"public_dns" yaml:"PublicDNS"`
	IPAddr     string `json:"ip_addr" yaml:"IPAddr"`

	// Auth
	SSLPem   string `json:"pem" yaml:"SSLPem"`
	SSLKey   string `json:"ssl_auth" yaml:"SSLKey"`
	RedisUri string `json:"redis_uri" yaml:"RedisUri"`
	RedisKey string `json:"redis_key" yaml:"RedisKey"`
}

var BaseConfig = initBase()

func initBase() *config {
	file, err := ioutil.ReadFile("./core/configs/config.yaml")
	if err != nil {
		if err := ioutil.WriteFile("./core/configs/config.yaml", []byte(cfp), 00666); err != nil {
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
ListenAddr: "0.0.0.0:8081"
Debug: true
PublicDNS: "8.8.8.8:53"
RedisUri: "127.0.0.1:6379"
RedisKey: ""
IPAddr: "0.0.0.0:8086"

SSLPem: "./configs/s1.pem"
SSLKey: "./configs/s1.key"
`
