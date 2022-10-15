package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
)

// RunContainerInitProcess 在容器内部执行，使用mount去挂载proc文件系统
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	// MS_NOEXEC是指本文件系统中不允许运行其他程序
	// MS_NOSUID是指本系统运行时不允许set userID或set groupID
	// MS_NODEV从Linux 2.4以来所有mount系统默认指定的参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}

	// 容器创建之后，执行的第一个进程是init进程而不是用户进程，这和预期是不一样的
	// 但是因为init进程的PID=1，如果kill掉则容器就退出了，因此在这里使用syscall.Exec
	// syscall.Exec最终调用了execve这个系统调用，执行当前filename对应的程序，并覆盖掉当前进程的镜像、数据及堆栈信息
	// 因此运行ps -ef时会发现/bin/bash是PID=1的进程
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
