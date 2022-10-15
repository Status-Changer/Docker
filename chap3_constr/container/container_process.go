//go:build linux

package container

import (
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
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
