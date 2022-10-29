package main

import (
	"Docker/cgroups"
	"Docker/cgroups/subsystems"
	"Docker/container"
	utils2 "Docker/utils"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// Run clone一个隔离namespace的进程，在子进程中调用自己发送init参数，调用init方法去初始化资源
func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig, volume string, containerName string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 记录容器的信息
	containerName, err := recordContainerInfo(parent.Process.Pid, commandArray, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// 创建Docker-cgroup，并设置对应的限制，随后初始化容器
	cgroupManager := cgroups.NewCgroupManager("Docker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArray, writePipe)
	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
	}

	mountURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkspace(rootURL, mountURL, volume)
	os.Exit(0)
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	defer writePipe.Close()
	command := strings.Join(commandArray, " ")
	log.Infof("all commands are %s", command)
	if _, err := writePipe.WriteString(command); err != nil {
		log.Errorf("writePipe.WriteString failed %v", err)
	}
}

// recordContainerInfo 记录容器的各种信息，并以json格式保存
func recordContainerInfo(containerPID int, commandArray []string, containerName string) (
	string, error) {
	id := utils2.RandStringBytes(utils2.DefaultContainerIdLength)
	createdTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	if len(containerName) == 0 {
		containerName = id
	}

	containerInfo := &container.ContainerInfo{
		Pid:         strconv.Itoa(containerPID),
		Id:          id,
		Name:        containerName,
		Command:     command,
		CreatedTime: createdTime,
		Status:      utils2.ContainerStatusRunning,
	}

	// 将容器信息序列化
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("record container info error. %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	// 容器存储的路径，如果不存在就级联创建
	dirUrl := fmt.Sprintf(utils2.DefaultContainerInfoStorageLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir all error %s. error %v", dirUrl, err)
		return "", err
	}

	// 创建最终的配置文件
	fileName := dirUrl + "/" + utils2.DefaultContainerInfoConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error. %v", fileName, err)
		return "", err
	}

	// 将序列化后的数据写入到文件
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(utils2.DefaultContainerInfoStorageLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}
