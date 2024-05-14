package api

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"log"
	"net"
	"strings"
	"time"
)

type HOST struct {
	Ip string `json:"ip"`
}

func DomainNameBinding(IpUrl string, url []string, hostname [][]string, token []string) {
	host := HOST{}
	body, err := HttpGet(IpUrl, nil)
	if err != nil {
		log.Println("获取公网IP失败")
	}

	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&host)
	IpAddr := host.Ip
	log.Println("外网ip：" + IpAddr)
	logs.Info("外网ip：" + IpAddr)
	checkUrl := hostname[0][0]
	ips, err := net.LookupIP(checkUrl)
	log.Println("checkUrl：" + checkUrl)
	logs.Info("checkUrl：" + checkUrl)
	urlIp := ""
	if err != nil {
		log.Println("无法查询到该checkUrl的IP地址")
		logs.Error("无法查询到该checkUrl的IP地址", err)
	} else {
		urlIp = ips[0].String()
		log.Println("checkUrl的IP:" + urlIp)
		logs.Info("checkUrl的IP:" + urlIp)
	}
	if IpAddr != urlIp {
		for i, u := range url {
			token := token[i]
			for _, hname := range hostname[i] {
				urlStr := ""
				var header = make(map[string]string)
				//判断字符串中是否包含字符
				if strings.Contains(u, "dynv6.com") {
					urlStr = u + "?hostname=" + hname + "&token=" + token + "&ipv4=" + IpAddr
				} else if strings.Contains(u, "dedyn.io") {
					urlStr = u + "?hostname=" + hname + "&myipv4=" + IpAddr
					header["Authorization"] = "Token " + token
				}
				if urlStr == "" {
					continue
				}
				log.Println(urlStr)
				response, err := HttpGet(urlStr, header)
				if err != nil {
					log.Println(err)
				}
				log.Println(hname + "绑定到" + IpAddr + ":" + string(response))
				logs.Info(hname + "绑定到" + IpAddr + ":" + string(response))
				time.Sleep(time.Second * 3)

			}
			time.Sleep(time.Second * 1)
		}

	} else {
		log.Println("外网ip与checkUrl的ip一致，无需更新")
		logs.Info("外网ip与checkUrl的ip一致，无需更新")
	}
}
