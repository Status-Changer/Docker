package main

import (
	"Docker/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func logContainer(containerName string) {
	// 找到日志文件的位置
	dirURL := fmt.Sprintf(utils.DefaultContainerInfoStorageLocation, containerName)
	logFileLocation := dirURL + utils.DefaultContainerLogFileName

	// 打开日志文件
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error. %v", logFileLocation, err)
		return
	}

	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error. %v", logFileLocation, err)
		return
	}
	fmt.Fprintf(os.Stdout, string(content))
}
