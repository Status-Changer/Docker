package main

import (
	"Docker/container"
	"Docker/utils"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// getContainerPidByName 根据容器名获得容器的Pid
func getContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(utils.DefaultContainerInfoStorageLocation, containerName)
	configFilePath := dirURL + containerName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.Info
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}

// ExecContainer
func ExecContainer(containerName string, commandArray []string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Exec container getContainerPidByName %s error. %v", containerName, err)
		return
	}

	// 把命令以空格为分隔符拼接成一个字符串
	commandString := strings.Join(commandArray, " ")

	// 为了再次运行nsenter.go中的C代码
	// 上述的C代码一旦启动就会运行，但是为了指定环境变量再运行一遍
	// 这里fork一个进程，并不在意各种namespace的隔离，带着环境变量去指定的namespace中进行操作
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(utils.EnvironmentExecPid, pid)
	os.Setenv(utils.EnvironmentExecCommand, commandString)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error. %v", containerName, err)
	}
}
