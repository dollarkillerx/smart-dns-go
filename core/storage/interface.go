package storage

import "golang.org/x/net/dns/dnsmessage"

type Storage interface {
	AnalysisDNS(domain string, dns *dnsmessage.Message, ip string) (result *dnsmessage.Message,err error)
}
