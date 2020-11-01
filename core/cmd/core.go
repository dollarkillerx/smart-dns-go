package main

import (
	"log"
	"net"

	"github.com/dollarkillerx/smart-dns-go/core/core"
	"github.com/dollarkillerx/smart-dns-go/core/pkg/config"
	"github.com/dollarkillerx/smart-dns-go/generate/gatewary"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if config.BaseConfig.Debug {
		log.SetFlags(log.Llongfile | log.LstdFlags)
	}

	log.Println("Listen Addr: ", config.BaseConfig.ListenAddr)

	listen, err := net.Listen("tcp", config.BaseConfig.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}

	creds, err := credentials.NewServerTLSFromFile(config.BaseConfig.SSLPem, config.BaseConfig.SSLKey)
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer(grpc.Creds(creds))
	core := core.New()
	gatewary.RegisterGatewaryServiceServer(server, core)
	if err := server.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
