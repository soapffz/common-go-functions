package pkg

import (
	"log"
	"os"
	"path/filepath"
)

func LoadFile(path string) []string {
	// 打开指定文件夹
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(0)
	}
	defer f.Close()
	// 读取目录下所有文件
	fileInfo, _ := f.ReadDir(-1)

	files := make([]string, 0)
	for _, info := range fileInfo {
		if filepath.Ext(info.Name()) == ".json" {
			files = append(files, filepath.Join(path+"/"+info.Name()))
		}
	}
	return files
}
