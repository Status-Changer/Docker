//go:build linux

package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 调用当前进程（/proc/self/exe）对创建出来的进程进行初始化：
//
// 1. 调用initCommand去初始化进程的环境和资源.
//
// 2. Cloneflags用于fork一个新进程，并使用namespace隔离外部环境.
//
// 3. 如果用户指定-ti参数就将当前进程的IO导入到标准IO上.
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// 传入管道文件读取端的句柄
	// 一个进程默认会有standard IO/error三个文件描述符，因此这里是第四个，将管道的一端传递给子进程
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = "/root/busybox"
	return cmd, writePipe
}

// NewPipe 生成两个匿名管道，用于读写
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
