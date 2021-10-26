package tools

import (
	"log"
	"testing"
)

func TestLong2String(t *testing.T) {

	appkey, _ := long2String(10000001)
	appid, _ := string2Long(appkey)
	checkapp := checkAppKey(appkey)
	log.Printf("TestLong2String.appid[%v] appkey[%v], checkapp[%v]", appid, appkey, checkapp)

	appkey1, appsecret1, _ := createAppkeyAndSecret(10000001, "terry", "WLAN")

	checkflag := CheckAppkeyAndSecret(appkey1, appsecret1)

	log.Printf("TestLong2String.appkey1[%v] appsecret1[%v], checkflag[%v]", appkey1, appsecret1, checkflag)
}

func TestGetNetworkName(t *testing.T) {
	testSerial := getCurrentSerial("WLAN")
	log.Printf("TestGetNetworkName.getCurrentSerial[%v]", testSerial)
}

func TestHTTPDns(t *testing.T) {
	// response, err := http.Get("http://203.107.1.33/174597/d?host=core.xinghuoyouxi.com")
	// if err != nil {
	// 	log.Printf("TestHTTPDns.err[%v]", err)
	// 	return
	// }

	// body, err1 := ioutil.ReadAll(response.Body)
	// if err1 != nil {
	// 	log.Printf("TestHTTPDns.err1[%v]", err1)
	// 	return
	// }

	// log.Printf("TestHTTPDns.body[%v]", string(body))
	ip := GetDnsIP("core.xinghuoyouxi.com")
	log.Printf("TestHTTPDns.ip[%v]", ip)

}
