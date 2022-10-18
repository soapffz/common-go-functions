package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 下载Github最新版本release中指定文件 by soapffz
// 原文章:https://blog.csdn.net/qq_39846820/article/details/115056151

func DownloadGithubRelease(repo string, localpath string, index string) []string {
	// repo:仓库名格式应为：360quake/quake_rs
	// localpath:保存文件的目录
	// index:release中的文件index号，类型为，0表示下载第一个文件，"1,3"表示下载第二和第四个文件

	r, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = r.Body.Close() }()
	body, _ := io.ReadAll(r.Body)

	var xxm mybody
	_ = json.Unmarshal(body, &xxm)
	indexArr := strings.Split(index, ",")
	filename_l := make([]string, 0)
	for _, i := range indexArr {
		//字符串转int
		i, _ := strconv.Atoi(i)
		url := xxm.Assets[i].BrowserDownloadURL
		filenameArr := strings.Split(url, "/")
		filename := filenameArr[len(filenameArr)-1]
		fmt.Println("正在下载---------" + filename)
		DownloadFileProgress(url, localpath+filename)
		filename_l = append(filename_l, filename)
	}
	return filename_l
}

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.Current += int64(n)
	fmt.Printf("\r进度 %.2f%%", float64(r.Current*10000/r.Total)/100)

	return
}

func DownloadFileProgress(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	reader := &Reader{
		Reader: r.Body,
		Total:  r.ContentLength,
	}

	_, _ = io.Copy(f, reader)
}

type mybody struct {
	Assets []Assets `json:"assets"`
}

type Assets struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}
