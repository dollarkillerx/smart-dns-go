package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	ListenAddr string `json:"listen_addr" yaml:"ListenAddr"`
	Debug      bool   `json:"debug" yaml:"Debug"`
}

var BaseConfig = initBase()

func initBase() *config {
	file, err := ioutil.ReadFile("./manager/configs/config.yaml")
	if err != nil {
		if err := ioutil.WriteFile("./manager/configs/config.yaml", []byte(cfp), 00666); err != nil {
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
ListenAddr: "0.0.0.0:8086"
Debug: true
`