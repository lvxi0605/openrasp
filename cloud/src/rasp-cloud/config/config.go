package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var TOMLConfig ServerConfig

// ServerInfo Config
type ServerConfig struct {
	AppKey        string                 `toml:"app_key"`        //app key
	AppSecret     string                 `toml:"app_secret"`     //app Secret
	AppVersion    string                 `toml:"app_version"`    //版本号
	AppMobile     string                 `toml:"app_mobile"`     //手机号
	AppEmail      string                 `toml:"app_email"`      //邮件号
	ServerNetwork string                 `toml:"server_network"` //服务器网卡名称
	AppId         string                 //服务器appid
	HTTPDNS       string                 `toml:"http_dns"`      //阿里http dns
	Location      LocationInfo           `toml:"location_info"` //本地经度
	Kafkaservers  map[string]KafkaServer `toml:"kafkaservers"`  //kafka服务器配置
	AttackTypes   map[string]AttackType  `toml:"attacktypes"`   //攻击类型
}

func (m ServerConfig) GetKafkaServer(name string) (KafkaServer, bool) {
	s, ok := m.Kafkaservers[name]
	return s, ok
}

func (m ServerConfig) GetAttackType(name string) (AttackType, bool) {
	s, ok := m.AttackTypes[name]
	return s, ok
}

//亚洲|中国|江西|南昌||电信|360100|China|CN|115.892151|28.676493
type LocationInfo struct {
	Continent     string `toml:"continent"`      //大陆
	Country       string `toml:"country"`        //国家
	Province      string `toml:"province"`       //省份
	City          string `toml:"city"`           //城市
	County        string `toml:"county"`         //区县
	ISP           string `toml:"isp"`            //运营商
	AreaCode      string `toml:"area_code"`      //区划编码
	CountryEN     string `toml:"country_en"`     //国家英文
	CountrySimple string `toml:"country_simple"` //国家简码
	Latitude      string `toml:"latitude"`       //经度
	Longitude     string `toml:"longitude"`      //纬度
}

/*
* KafkaServer服务器配置
 */
type KafkaServer struct {
	Addr     string `toml:"KafkaAddr"`
	User     string `toml:"KafkaUser"`
	Password string `toml:"KafkaPwd"`
	Enable   bool   `toml:"KafkaEnable"`
	Topic    string `toml:"KafkaTopic"`
	Port     string `toml:"KafkaPort"`
}

type AttackType struct {
	ID    int    `toml:"AttackID"`
	Type  string `toml:"AttackType"`
	Name  string `toml:"AttackName"`
	Level string `toml:"AttackLevel"`
}

func init() {

	TOMLConfig.HTTPDNS = "http://203.107.1.33/174597/d"

	if _, err := toml.DecodeFile("conf/local.toml", &TOMLConfig); err != nil {
		return
	}

	log.Printf("DEBUG: Config: %#v\n", TOMLConfig)
}
