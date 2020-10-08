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
	IPS_LT  IPS = "联通"
	IPS_DX  IPS = "电信"
	IPS_YD  IPS = "移动"
	IPS_JY  IPS = "教育网"
	IPS_TT  IPS = "铁通"
	IPS_PBS IPS = "鹏博士"
	IPS_QT  IPS = "0"
)

type CITY string

const (
	CITY_GN CITY = "中国"
	CITY_QT CITY = "0"
)
