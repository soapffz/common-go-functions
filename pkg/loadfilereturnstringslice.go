package pkg

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// 传入文件名，按行读取返回数组列表 by soapffz 2022-10-23

func ReadFileReturnStringSlice(filename string) []string {
	fileIn, fileInErr := os.Open(filename)
	if fileInErr != nil {
		log.Fatal("[Warn] 打开文件失败：", filename)
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	var fileList []string
	for {
		inputString, readerError := finReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		inputString = strings.Replace(inputString, "\n", "", -1) // 去掉回车
		fileList = append(fileList, inputString)
	}
	return fileList
}
