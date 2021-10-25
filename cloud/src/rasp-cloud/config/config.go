package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var TOMLConfig ServerConfig

// ServerInfo Config
type ServerConfig struct {
	AppKey        string `toml:"app_key"`        //app key
	AppSecret     string `toml:"app_secret"`     //app Secret
	AppVersion    string `toml:"app_version"`    //版本号
	AppMobile     string `toml:"app_mobile"`     //手机号
	AppEmail      string `toml:"app_email"`      //邮件号
	ServerNetwork string `toml:"server_network"` //服务器网卡名称
}

/*
* KafkaServer服务器配置
 */
type KafkaServer struct {
	Addr     string `toml:"addr"`
	User     string `toml:"addr"`
	Password string `toml:"password"`
	Enable   bool
	DB       int `toml:"db"`
}

func init() {

	if _, err := toml.DecodeFile("conf/local.toml", &TOMLConfig); err != nil {
		return
	}
	log.Printf("DEBUG: Config: %#v\n", TOMLConfig)
}
