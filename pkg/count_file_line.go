package pkg

import (
	"bufio"
	"os"
)

func CountFileLine(filename string) (line_num int) {
	//传入文件路径，返回文件行数
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	return count
}
