package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
)

// RunContainerInitProcess 在容器内部执行，使用mount去挂载proc文件系统
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
