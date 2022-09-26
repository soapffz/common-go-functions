package pkg

import (
	"log"
	"net/http"
	"strings"
)

func PushMsgByServerJ(serverj_key string, title string, content string) {
	// 传入server酱的key，title和内容使用默认推送方式进行推送
	url := "http://sc.ftqq.com/" + serverj_key + ".send?text=" + title
	_, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("&desp="+content))
	if err != nil {
		log.Fatalln("推送失败，请检查配置或网络")
	} else {
		// fmt.Println(resp)
		log.Println("推送消息成功")
	}
}
