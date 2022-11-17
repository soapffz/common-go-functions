package pkg

import (
	"bufio"
	"log"
	"os"
)

func ReadSpecifiedLineInFile(filepath string, begin_line int, end_line int) []string {
	// 读取指定文件指定的行，从指定开始行读到指定结束行，若开始行大于文件行数则报错，若指定结束行大于文件行数则读到结尾返回

	if begin_line < 1 || begin_line > end_line {
		log.Fatal("输入参数不对，程序退出")
		os.Exit(-1)
	}
	_, err := os.Stat(filepath)
	if err == nil {
		// 文件存在则读取文件
		var content []string
		filehandle, err := os.Open(filepath)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		defer filehandle.Close()
		linecount := 1
		fileScanner := bufio.NewScanner(filehandle)
		for fileScanner.Scan() {
			if linecount >= begin_line && linecount <= end_line {
				content = append(content, fileScanner.Text()+"\n")
			}
			linecount++
		}
		return content
	}
	if os.IsNotExist(err) {
		log.Fatal("文件不存在，请检查文件是否存在")
		os.Exit(-1)
	}
	return nil
}
