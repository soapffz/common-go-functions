package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type IpDomain struct {
	Domain string `json:"domain"`
	Title  string `json:"title"`
}

func Ip2DomainAndWebWeight(ip string) (string, int) {
	// 根据传入的ip地址解析域名及爱站中的网站权重
	// 参考项目https://github.com/Sma11New/ip2domain
	ip = strings.ReplaceAll(ip, " ", "")
	domain := Ip2Domain(ip)
	if domain != "" {
		web_weight := AiZhanRankQuery(domain)
		return domain, web_weight
	}
	return "", 0
}

func Ip2Domain(ip string) string {
	// 传入ip，解析返回最短域名
	ip = strings.ReplaceAll(ip, " ", "")
	client := &http.Client{}
	url := "https://api.webscan.cc/?action=query&ip=" + ip
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36")
	response, _ := client.Do(reqest)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		var ll []IpDomain = make([]IpDomain, 0)
		res_byte, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(res_byte, &ll)
		if err != nil {
			return ""
		}
		if len(ll) > 1 {
			// 存在多个数据时，取最短长度值
			var min_domain string
			for _, data := range ll {
				domain := data.Domain
				if min_domain == "" {
					min_domain = domain
				} else if len(domain) < len(min_domain) {
					min_domain = domain
				} else {
					continue
				}
			}
			return min_domain
		} else if len(ll) == 0 {
			// 数组只有一个值，直接取值返回
			return ll[0].Domain
		} else {
			// 没有获取到域名
			return ""
		}
	}
	return ""
}

func AiZhanRankQuery(domain string) int {
	// 传入域名，返回爱站查询到的权重数据
	// 正常返回权重，查询失败或其他情况返回0
	client := &http.Client{}
	url := "https://www.aizhan.com/cha/" + domain + "/"
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Host", "www.aizhan.com")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:61.0) Gecko/20100101 Firefox/61.0")
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, _ := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		res, _ := ioutil.ReadAll(response.Body)
		resstring := string(res)
		reg := regexp.MustCompile("aizhan.com/images/br/(.*?).png").FindStringSubmatch(resstring)
		if len(reg) == 2 {
			str_web_weight := reg[1]
			// fmt.Println(domain + " 的权重是：" + web_weight)
			web_weight, _ := strconv.Atoi(str_web_weight)
			return web_weight
		} else {
			return 0
		}
	}
	return 0
}
