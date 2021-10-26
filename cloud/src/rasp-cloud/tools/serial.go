package tools

import (
	"crypto/sha256"
	"fmt"
	"net"
	"runtime"

	"golang.org/x/crypto/pbkdf2"
)

//获取当前机器的serial
func getCurrentSerial(networkName string) string {
	if networkName == "" {
		if runtime.GOOS == "linux" {
			networkName = "eth0"
		} else if runtime.GOOS == "darwin" {
			networkName = "en0"
		} else if runtime.GOOS == "windows" {
			networkName = "WLAN"
		}
	}

	//获取所有网卡
	a, err := net.Interfaces()
	if err != nil {
		return ""
	}

	key := "FwecnH3ibu"

	//查找eth0的网卡
	// var ipString, hardwareAddr string
	var hardwareAddr string

	for i := 0; i < len(a); i++ {
		if a[i].Name == networkName {
			addr, err1 := a[i].Addrs()
			if err1 != nil {
				fmt.Println(err1)
				continue
			}

			if len(addr) <= 0 {
				// fmt.Println("len addr is %v", len(addr))
				continue
			}

			for _, address := range addr {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						hardwareAddr = fmt.Sprintf("%v", a[i].HardwareAddr)
						// ipString = ipnet.IP.String()
						break
					}
				}
			}
			break
		}
	}

	// if ipString == "" || hardwareAddr == "" {
	// 	fmt.Println("get eth0 error")
	// 	return ""
	// }
	if hardwareAddr == "" {
		fmt.Println("get eth0 error")

		return ""
	}

	// log.Printf("getCurrentSerial. hardwareAddr[%v], runtime.GOOS[%v]", hardwareAddr, runtime.GOOS)

	// passwordStr := fmt.Sprintf("%v%v", ipString, hardwareAddr)
	passwordStr := fmt.Sprintf("%v", hardwareAddr)
	// passwordStr = "f0:d4:e2:e8:89:bc"
	newPasswd := pbkdf2.Key([]byte(passwordStr), []byte(key), 10000, 50, sha256.New)
	Passwd := fmt.Sprintf("%x", newPasswd)
	return Passwd
}
