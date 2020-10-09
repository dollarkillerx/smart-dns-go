package smart_dns_go

import (
	"fmt"
	"log"
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
