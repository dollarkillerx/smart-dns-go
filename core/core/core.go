package core

import (
	"context"
	"github.com/dollarkillerx/smart-dns-go/generate/gatewary"
)

type Core struct{}

func (c *Core) DNSLookup(ctx context.Context, req *gatewary.DnsRequest) (resp *gatewary.DnsResponse, err error) {

}
