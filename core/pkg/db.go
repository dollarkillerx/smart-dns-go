package pkg

// 结构定义
type DNSNode struct {
	Country string
	ISP     string

	Value string
	TTL   uint32
	Pref  uint16
}

// domain: type: id
// id: location: value

type IPS string

const (
	IpsLt      IPS = "联通"
	IpsDx      IPS = "电信"
	IpsYd      IPS = "移动"
	IpsJy      IPS = "教育网"
	IpsTt      IPS = "铁通"
	IpsPbs     IPS = "鹏博士"
	IpsDefault IPS = "0"
)

func (i IPS) String() string {
	return string(i)
}

type CITY string

const (
	CityGn      CITY = "中国"
	CityDefault CITY = "0"
)

func (i CITY) String() string {
	return string(i)
}
