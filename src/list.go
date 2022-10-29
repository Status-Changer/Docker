package main

import (
	"Docker/src/container"
	"Docker/src/utils"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

// ListAllContainers 展示所有的容器
func ListAllContainers() {
	dirURL := fmt.Sprintf(utils.DefaultContainerInfoStorageLocation, "")
	dirURL = dirURL[:len(dirURL)-1] // 去掉一个末尾的'/'，此时是作为目录

	// 读取该目录下的所有文件
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		log.Errorf("Read dir %s error. %v", dirURL, err)
		return
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		tempContainer, err := getContainerInfo(file)
		if err != nil {
			log.Errorf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tempContainer)
	}

	// tabwriter用于在控制台打印对齐的表格
	writer := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(writer, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id, item.Name, item.Pid, item.Status, item.Command, item.CreatedTime)
	}
	// 刷新stdout缓冲区，使之立刻打印
	if err := writer.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
	// 根据文件名及相关的路径配置信息，获得文件绝对路径
	containerName := file.Name()
	configFileDir := fmt.Sprintf(utils.DefaultContainerInfoStorageLocation, containerName)
	configFileDir = configFileDir + utils.DefaultContainerInfoConfigName

	// 读取json信息并反序列化成容器信息对象
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("Read file %s error %v", configFileDir, err)
		return nil, err
	}

	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
