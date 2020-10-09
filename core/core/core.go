package core

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"

	"github.com/dollarkillerx/smart-dns-go/core/storage"
	"github.com/dollarkillerx/smart-dns-go/core/storage/simple"
	"github.com/dollarkillerx/smart-dns-go/generate/gatewary"
	"golang.org/x/net/dns/dnsmessage"
)

type Core struct {
	Db storage.Storage
}

func New() *Core {
	return &Core{
		Db: simple.NewSimple(),
	}
}

func (c *Core) DNSLookup(ctx context.Context, req *gatewary.DnsRequest) (resp *gatewary.DnsResponse, err error) {
	var m dnsmessage.Message
	if err := m.Unpack(req.Message); err != nil {
		log.Println(err)
		return nil, err
	}

	dns, err := c.Db.AnalysisDNS(m.Questions[0].Name.String(), &m, req.Ip)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, err
		}
		log.Println(err)
		return nil, err
	}

	pack, err := dns.Pack()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &gatewary.DnsResponse{
		Message: pack,
	}, nil
}
