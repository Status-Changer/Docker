package main

import (
	"Docker/chap3_constr/container"
	log "github.com/sirupsen/logrus"
	"os"
)

// Run clone一个隔离namespace的进程，在子进程中调用自己发送init参数，调用init方法去初始化资源
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
