package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	parser "github.com/Cgboal/DomainParser"
)

var extractor parser.Parser

type IpDomain struct {
	Domain string `json:"domain"`
	Title  string `json:"title"`
}

type TooLTT struct {
	Status int  `json:"status"`
	Data   Data `json:"data"`
}
type InData struct {
	Domain string `json:"domain"`
	Value  string `json:"value"`
	Type   string `json:"type"`
	Time   string `json:"time"`
}
type Data struct {
	Data  []InData `json:"data"`
	Count int      `json:"count"`
	Type  int      `json:"type"`
}

func Ip2DomainAndWebWeight(ip string) (string, string, int) {
	// 传入ip地址，返回域名，根域名，根域名权重
	// 参考项目https://github.com/Sma11New/ip2domain
	// https://github.com/Cgboal/DomainParser
	// https://github.com/MoYang233/subdomain-demo/blob/main/subdomain-demo2/util/util.go

	ip = strings.ReplaceAll(ip, " ", "")
	// 使用ip解析的主函数获得域名列表
	domain_l := Ip2Domain(ip)
	if len(domain_l) != 0 {
		BIG_WEIGHT := 0
		BIG_WEIGHT_DOMAIN := ""
		BIG_WEIGHT_ROOT_DOMAIN := ""
		for _, domain := range domain_l {
			// 使用域名列表获得根域名列表
			root_domain := GetRootDomain(domain)
			if root_domain != "" {
				web_weight := AiZhanRankQuery(root_domain)
				if web_weight >= BIG_WEIGHT {
					BIG_WEIGHT = web_weight
					BIG_WEIGHT_DOMAIN = domain
					BIG_WEIGHT_ROOT_DOMAIN = root_domain
				}
			}
			// 根据根域名列表查询权重，选择一个权重最高的，若相同随机取一个
		}
		return BIG_WEIGHT_DOMAIN, BIG_WEIGHT_ROOT_DOMAIN, BIG_WEIGHT
	}
	return "", "", 0
}

func GetRootDomain(domain string) string {
	// 传入域名，解析为根域名
	if domain != "" {
		// 提取顶级域名
		extractor = parser.NewDomainParser()
		return extractor.GetFQDN(domain)
	}
	return ""
}

func Ip2Domain(ip string) []string {
	// ip解析为域名的主函数，传入一个ip，经过多种方法处理后，返回最小的域名列表

	var domain_l []string
	ip = strings.ReplaceAll(ip, " ", "")
	domain_by_web_scan := Ip2DomainByWebScancc(ip)
	domain_by_dns_grep := Ip2DomainByDnsGrep(ip)
	domain1 := CleanDomain(domain_by_web_scan)
	domain2 := CleanDomain(domain_by_dns_grep)
	for _, i := range []string{domain1, domain2} {
		if i != "" {
			domain_l = append(domain_l, i)
		}
	}
	return domain_l
}

func CleanDomain(domain string) string {
	// 传入domain，清理
	// 去除前后空格
	domain = strings.TrimSpace(domain)
	// 如果是domain/ip:port模式，去除掉端口
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}
	// 去除www前缀
	if strings.HasPrefix(domain, "www.") {
		domain = strings.Trim(domain, "www.")
	}
	// 如果传入的是个ip，返回空字符串
	match, _ := regexp.MatchString("[A-Za-z]", domain)
	if !match {
		return ""
	}
	if domain == "" {
		return ""
	} else {
		return domain
	}
}

func Ip2DomainByWebScancc(ip string) string {
	// 传入ip，根据api.webscan.cc接口解析返回最短域名
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
		} else if len(ll) == 1 {
			// 数组只有一个值，直接取值返回
			return ll[0].Domain
		} else {
			// 没有获取到域名
			return ""
		}
	}
	return ""
}

func Ip2DomainByDnsGrep(ip string) string {
	// 传入ip，根据dnsgrep.cn接口解析返回最短域名
	ip = strings.ReplaceAll(ip, " ", "")
	client := &http.Client{}
	url := "https://www.dnsgrep.cn/api/query?q=" + ip + "&token=6fecc6d76090e8fd4ff0ebaa9af30c7d"
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36")
	response, _ := client.Do(reqest)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		var dd TooLTT
		res_byte, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(res_byte, &dd)
		if err != nil {
			return ""
		}
		if dd.Data.Count == 1 {
			return dd.Data.Data[0].Domain
		} else if dd.Data.Count > 1 {
			var min_domain string
			for _, data := range dd.Data.Data {
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
		} else {
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
