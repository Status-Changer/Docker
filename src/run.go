package main

import (
	"Docker/src/cgroups"
	"Docker/src/cgroups/subsystems"
	"Docker/src/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// Run clone一个隔离namespace的进程，在子进程中调用自己发送init参数，调用init方法去初始化资源
func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 创建Docker-cgroup，并设置对应的限制，随后初始化容器
	cgroupManager := cgroups.NewCgroupManager("Docker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(commandArray, writePipe)

	parent.Wait()
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	defer writePipe.Close()
	command := strings.Join(commandArray, " ")
	log.Infof("all commands are %s", command)
	if _, err := writePipe.WriteString(command); err != nil {
		log.Errorf("writePipe.WriteString failed %v", err)
	}
}
