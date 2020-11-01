package test

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestP(t *testing.T) {
	ip := "0.0.0.0:86858"

	c := strings.Index(ip, ":")
	if c != -1 {
		ip = ip[:c]
	}
	log.Println(ip)
}

func TestPi(t *testing.T) {
	pool := newRedisPool()
	red := pool.Get()
	defer red.Close()

	if _, err := red.Do("HDEL", "cpx", "sk"); err != nil {
		log.Fatalln(err)
	}

	do, err := redis.String(red.Do("HGET", "cpx", "sk"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(do)
}

type DNSNode struct {
	Country string
	ISP     string

	Value string
	TTL   int
}

func TestCpx(t *testing.T) {
	var tag DNSNode
	log.Println(tag)
}

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,                // 池中的最大空闲连接数
		MaxActive:   30,                // 最大连接数
		IdleTimeout: 300 * time.Second, // 超时回收
		Dial: func() (conn redis.Conn, e error) {
			// 1. 打开连接
			dial, e := redis.Dial("tcp", "0.0.0.0:6379")
			if e != nil {
				log.Fatalln(e)
			}
			// 2. 访问认证
			//if config.BaseConfig.RedisKey != "" {
			//	dial.Do("AUTH", config.BaseConfig.RedisKey)
			//}
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

func TestJm(t *testing.T) {
	p := "135.4.5.6"
	log.Println(ParseIP(p))
}

func ParseIP(p string) ([4]byte, error) {
	split := strings.Split(p, ".")
	if len(split) != 4 {
		return [4]byte{}, fmt.Errorf("not ep")
	}
	cp := [4]byte{}
	for i, v := range split {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			return [4]byte{}, err
		}
		cp[i] = uint8(atoi)
	}

	return cp, nil
}

func TestAbc(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("dig @0.0.0.0 -p 8986 www.baidu.com")
	listen, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8986})
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()

	for {
		buf := make([]byte, 512)
		i, addr, err := listen.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go dnsCore(buf[:i], addr, listen)
	}
}

func dnsCore(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("Recover Err: ", e)
			return
		}
	}()

	var msg dnsmessage.Message
	if err := msg.Unpack(data); err != nil {
		log.Println(err)
		return
	}
	log.Println(msg)
	marshal, err := json.Marshal(msg)
	if err == nil {
		log.Println(string(marshal))
	}
	pack, err := msg.Pack()
	if err != nil {
		log.Println(err)
		return
	}

	dns, err := dialDns(pack)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := conn.WriteToUDP(dns, addr); err != nil {
		log.Println(err)
		return
	}
}

func dialDns(msg []byte) ([]byte, error) {
	conn, err := net.DialTimeout("udp", "223.5.5.5:53", time.Second)
	if err != nil {
		return nil, err
	}

	conn.SetWriteDeadline(time.Now().Add(time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Second))
	conn.SetDeadline(time.Now().Add(time.Second))

	if _, err := conn.Write(msg); err != nil {
		return nil, err
	}

	buf := make([]byte, 512)
	var m dnsmessage.Message
	read, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if err := m.Unpack(buf[:read]); err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(m)
	if err == nil {
		log.Println(string(marshal))
	}

	pack, err := m.Pack()
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func dialDns2(msg []byte) (*dnsmessage.Message, error) {
	conn, err := net.DialTimeout("udp", "223.5.5.5:53", time.Second)
	if err != nil {
		return nil, err
	}

	conn.SetWriteDeadline(time.Now().Add(time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Second))
	conn.SetDeadline(time.Now().Add(time.Second))

	if _, err := conn.Write(msg); err != nil {
		return nil, err
	}

	buf := make([]byte, 512)
	var m dnsmessage.Message
	read, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if err := m.Unpack(buf[:read]); err != nil {
		return nil, err
	}

	return &m, nil
}

func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func TestP2(t *testing.T) {
	dnsAis("baidu.com.")
}

func dnsAis(domain string) {
	var m dnsmessage.Message
	m.Header.ID = uint16(Random(1000, 65534))
	//m.Header.RecursionDesired = true
	m.Questions = append(m.Questions, dnsmessage.Question{
		Name:  dnsmessage.MustNewName(domain),
		Type:  dnsmessage.TypeA,
		Class: dnsmessage.ClassINET,
	})

	pack, err := m.Pack()
	if err != nil {
		log.Fatalln(err)
	}

	rm, err := dialDns2(pack)
	if err != nil {
		log.Println(err)
	}

	log.Println(rm.Answers[0].Body.GoString())
}
