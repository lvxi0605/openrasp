package tools

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rasp-cloud/config"
)

func GetDnsIP(host string) string {
	type HTTPDnsRes struct {
		Host      string   `json:"host"`
		IPS       []string `json:"ips"`
		TTL       int      `json:"ttl"`
		OriginTTL int      `json:"origin_ttl"`
		ClientIP  string   `json:"client_ip"`
	}

	httpHost := config.TOMLConfig.HTTPDNS + "?host=" + host
	response, err := http.Get(httpHost)
	if err != nil {
		log.Printf("GetDnsIP.err[%v]", err)
		return host
	}

	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		log.Printf("GetDnsIP.err1[%v]", err1)
		return host
	}

	var httpres HTTPDnsRes

	err2 := json.Unmarshal(body, &httpres)
	if err2 != nil {
		return host
	}

	return httpres.IPS[0]

}
