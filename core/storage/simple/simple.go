package simple

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dollarkillerx/smart-dns-go/core/pkg"
	"github.com/dollarkillerx/smart-dns-go/core/pkg/config"
	"github.com/dollarkillerx/smart-dns-ip/generate/smart_dns_ip"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/net/dns/dnsmessage"
	"google.golang.org/grpc"
)

type simple struct {
	redisConn *redis.Pool
	ipSearch  smart_dns_ip.IPSearchClient
}

func NewSimple() *simple {
	dial, err := grpc.Dial(config.BaseConfig.IPAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return &simple{
		redisConn: newRedisPool(),
		ipSearch:  smart_dns_ip.NewIPSearchClient(dial),
	}
}

func (s *simple) AnalysisDNS(domain string, dns *dnsmessage.Message, ip string) (result *dnsmessage.Message, err error) {
	conn := s.redisConn.Get()
	defer conn.Close()
	if domain == "" || ip == "" {
		return nil, fmt.Errorf("domain is null")
	}
	c := strings.Index(ip, ":")
	if c != -1 {
		ip = ip[:c]
	}

	do, err := redis.String(conn.Do("HGET", domain, dns.Questions[0].Type.String()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ls []pkg.DNSNode
	if err := json.Unmarshal([]byte(do), &ls); err != nil {
		log.Println(err)
		return nil, err
	}

	// check ip
	search, err := s.ipSearch.IPSearch(context.TODO(), &smart_dns_ip.IPSearchRequest{
		Ip: ip,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// sp
	country := ""
	if search.Country == "0" {
		country = "Def"
	} else {
		country = search.Country
	}
	isp := ""
	if search.Isp == "0" {
		isp = "Def"
	} else {
		isp = search.Isp
	}
	// check
	var def []pkg.DNSNode
	var tag []pkg.DNSNode
	for _, v := range ls {
		if v.Country == "Def" && v.ISP == "Def" {
			def = append(def, v)
			continue
		}

		if v.Country == country && v.ISP == isp {
			tag = append(tag, v)
		}
	}
	if len(def) == 0 {
		return nil, fmt.Errorf("not data")
	}
	if len(tag) == 0 {
		def = tag
	}

	result = dns
	result.Response = true
	result.RecursionAvailable = true
	result.Additionals = []dnsmessage.Resource{}
	result.Answers = []dnsmessage.Resource{}
	switch dns.Questions[0].Type {
	case dnsmessage.TypeA:
		for _, v := range def {
			result.Answers = append(result.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(domain),
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
					TTL:   v.TTL,
				},
				Body: &dnsmessage.AResource{A: [4]byte{46, 82, 174, 69}},
			})
		}
	case dnsmessage.TypeAAAA:

	case dnsmessage.TypeCNAME:
		for _, v := range def {
			result.Answers = append(result.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(domain),
					Type:  dnsmessage.TypeCNAME,
					Class: dnsmessage.ClassINET,
					TTL:   v.TTL,
				},
				Body: &dnsmessage.CNAMEResource{CNAME: dnsmessage.MustNewName(v.Value)},
			})
		}
	case dnsmessage.TypeTXT:
		for _, v := range def {
			result.Answers = append(result.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(domain),
					Type:  dnsmessage.TypeTXT,
					Class: dnsmessage.ClassINET,
					TTL:   v.TTL,
				},
				Body: &dnsmessage.TXTResource{TXT: []string{v.Value}},
			})
		}
	case dnsmessage.TypeMX:
		for _, v := range def {
			result.Answers = append(result.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(domain),
					Type:  dnsmessage.TypeMX,
					Class: dnsmessage.ClassINET,
					TTL:   v.TTL,
				},
				Body: &dnsmessage.MXResource{Pref: v.Pref, MX: dnsmessage.MustNewName(v.Value)},
			})
		}
	}
}

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,                // 池中的最大空闲连接数
		MaxActive:   30,                // 最大连接数
		IdleTimeout: 300 * time.Second, // 超时回收
		Dial: func() (conn redis.Conn, e error) {
			// 1. 打开连接
			dial, e := redis.Dial("tcp", config.BaseConfig.RedisUri)
			if e != nil {
				log.Fatalln(e)
			}
			// 2. 访问认证
			if config.BaseConfig.RedisKey != "" {
				dial.Do("AUTH", config.BaseConfig.RedisKey)
			}
			return dial, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 定时检查连接是否可用
			// time.Since(t) 获取离现在过了多少时间
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
