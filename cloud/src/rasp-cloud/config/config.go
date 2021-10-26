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
	HTTPDNS       string                 `toml:"http_dns"`     //阿里http dns
	Kafkaservers  map[string]KafkaServer `toml:"kafkaservers"` //kafka服务器配置
}

func (m ServerConfig) GetKafkaServer(name string) (KafkaServer, bool) {
	s, ok := m.Kafkaservers[name]
	return s, ok
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

func init() {

	TOMLConfig.HTTPDNS = "http://203.107.1.33/174597/d"

	if _, err := toml.DecodeFile("conf/local.toml", &TOMLConfig); err != nil {
		return
	}

	log.Printf("DEBUG: Config: %#v\n", TOMLConfig)
}
