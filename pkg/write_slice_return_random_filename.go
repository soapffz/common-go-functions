package pkg

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// 在当前文件将string数组写入随机文件，返回文件名

func WriteSliceReturnRandomFilename(data_l []string) (filename string) {
	// 传入字符串
	filePath := "./tmp_" + time.Now().Format("2006_01_02_15_04_05") + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("资产文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for _, singledata := range data_l {
		write.WriteString(singledata + "\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return filePath
}
