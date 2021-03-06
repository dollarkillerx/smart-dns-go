package core

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/dollarkillerx/smart-dns-go/gatewary/pkg/config"
	"github.com/dollarkillerx/smart-dns-go/generate/gatewary"
	"golang.org/x/net/dns/dnsmessage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Core struct {
	core gatewary.GatewaryServiceClient
}

func NewCore() *Core {
	var opts []grpc.DialOption
	creds, err := credentials.NewClientTLSFromFile(config.BaseConfig.SSLPem, "smart_dns_core")
	if err != nil {
		log.Fatalln(err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	//opts = append(opts, grpc.WithPerRPCCredentials(&customCredential{}))
	dial, err := grpc.Dial(config.BaseConfig.CoreAddr, opts...)
	if err != nil {
		log.Fatalln(err)
	}

	return &Core{
		core: gatewary.NewGatewaryServiceClient(dial),
	}
}

//// customCredential自定义认证
//type customCredential struct{}
//func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
//	return map[string]string{
//		"appid":  "0001",
//		"appkey": "key",
//	}, nil
//}
//func (c *customCredential) RequireTransportSecurity() bool {
//	return true
//}

func (c *Core) Core() error {
	addr, err := net.ResolveUDPAddr("udp", config.BaseConfig.ListenAddr)
	if err != nil {
		log.Fatalln("Can't resolve address: ", err)
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()

	for {
		buf := make([]byte, 512)
		i, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		go c.coreDnsServer(udpConn, addr, buf[:i])
	}
}

func (c *Core) coreDnsServer(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	var msg dnsmessage.Message
	if err := msg.Unpack(data); err != nil {
		log.Println(string(data))
		log.Println(err)
		return
	}

	//log.Println(msg.GoString())
	if len(msg.Questions) != 1 {
		dns, err := c.publicDNS(data)
		if err != nil {
			log.Println(err)
			return
		}

		c.respDNS(conn, addr, dns)
		return
	}

	dnsMsg, err := msg.Pack()
	if err != nil {
		log.Println(err)
		return
	}

	//var bm dnsmessage.Message
	switch msg.Questions[0].Type {
	case dnsmessage.TypeA:
		if err := c.PSend(dnsMsg, conn, addr); err == nil {
			return
		}
	case dnsmessage.TypeCNAME:
		if err := c.PSend(dnsMsg, conn, addr); err == nil {
			return
		}
	case dnsmessage.TypeMX:
		if err := c.PSend(dnsMsg, conn, addr); err == nil {
			return
		}
	case dnsmessage.TypeTXT:
		if err := c.PSend(dnsMsg, conn, addr); err == nil {
			return
		}
	}

	dns, err := c.publicDNS(data)
	if err != nil {
		log.Println(err)
		return
	}

	c.respDNS(conn, addr, dns)
}

func (c *Core) publicDNS(msg []byte) (*dnsmessage.Message, error) {
	dial, err := net.Dial("udp", config.BaseConfig.PublicDNS)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dial.SetWriteDeadline(time.Now().Add(time.Second * 3))
	dial.SetReadDeadline(time.Now().Add(time.Second * 3))
	dial.SetDeadline(time.Now().Add(time.Second * 3))
	_, err = dial.Write(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	buf := make([]byte, 512)
	var m dnsmessage.Message
	read, err := dial.Read(buf)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := m.Unpack(buf[:read]); err != nil {
		log.Println(err)
		return nil, err
	}

	return &m, nil
}

func (c *Core) respDNS(conn *net.UDPConn, addr *net.UDPAddr, data *dnsmessage.Message) {
	pack, err := data.Pack()
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(pack, addr); err != nil {
		log.Println(err)
	}
}

func (c *Core) PSend(dnsMsg []byte, conn *net.UDPConn, addr *net.UDPAddr) error {
	lookup, err := c.core.DNSLookup(context.TODO(), &gatewary.DnsRequest{
		Message: dnsMsg,
		Ip:      addr.String(),
	})
	if err != nil {
		if strings.Index(err.Error(),"nil returned") != -1 {
			return err
		}
		log.Println(err)
		return err
	}
	c.respDNSMsg(conn,addr,lookup.Message)
	return nil
}

func (c *Core) respDNSMsg(conn *net.UDPConn, addr *net.UDPAddr, pack []byte) {
	if _, err := conn.WriteToUDP(pack, addr); err != nil {
		log.Println(err)
	}
}

// dns.Questions[0].Name.String() www.google.com.
// mx    {"ID":29208,"Response":true,"OpCode":0,"Authoritative":false,"Truncated":false,"RecursionDesired":true,"RecursionAvailable":true,"RCode":0,"Questions":[{"Name":{"Data":[100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":17},"Type":15,"Class":1}],"Answers":[{"Header":{"Name":{"Data":[100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":17},"Type":15,"Class":1,"TTL":300,"Length":15},"Body":{"Pref":10,"MX":{"Data":[101,109,120,46,109,97,105,108,46,114,117,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":12}}}],"Authorities":[],"Additionals":[]}
// a     {"ID":8971,"Response":true,"OpCode":0,"Authoritative":false,"Truncated":false,"RecursionDesired":true,"RecursionAvailable":true,"RCode":0,"Questions":[{"Name":{"Data":[119,119,119,46,100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":21},"Type":1,"Class":1}],"Answers":[{"Header":{"Name":{"Data":[119,119,119,46,100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":21},"Type":1,"Class":1,"TTL":300,"Length":4},"Body":{"A":[104,27,129,116]}},{"Header":{"Name":{"Data":[119,119,119,46,100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":21},"Type":1,"Class":1,"TTL":300,"Length":4},"Body":{"A":[172,67,197,133]}},{"Header":{"Name":{"Data":[119,119,119,46,100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":21},"Type":1,"Class":1,"TTL":300,"Length":4},"Body":{"A":[104,27,128,116]}}],"Authorities":[],"Additionals":[]}
// cname {"ID":47799,"Response":true,"OpCode":0,"Authoritative":false,"Truncated":false,"RecursionDesired":true,"RecursionAvailable":true,"RCode":0,"Questions":[{"Name":{"Data":[99,112,120,46,100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":21},"Type":5,"Class":1}],"Answers":[],"Authorities":[{"Header":{"Name":{"Data":[100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":17},"Type":6,"Class":1,"TTL":3600,"Length":46},"Body":{"NS":{"Data":[97,109,121,46,110,115,46,99,108,111,117,100,102,108,97,114,101,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":22},"MBox":{"Data":[100,110,115,46,99,108,111,117,100,102,108,97,114,101,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":19},"Serial":2035320884,"Refresh":10000,"Retry":2400,"Expire":604800,"MinTTL":3600}}],"Additionals":[]}
// txt   {"ID":59905,"Response":true,"OpCode":0,"Authoritative":false,"Truncated":false,"RecursionDesired":true,"RecursionAvailable":true,"RCode":0,"Questions":[{"Name":{"Data":[100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":17},"Type":16,"Class":1}],"Answers":[{"Header":{"Name":{"Data":[100,111,108,108,97,114,107,105,108,108,101,114,46,99,111,109,46,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"Length":17},"Type":16,"Class":1,"TTL":300,"Length":29},"Body":{"TXT":["v=spf1 redirect=_spf.mail.ru"]}}],"Authorities":[],"Additionals":[]}
