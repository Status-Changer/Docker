package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cGroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 512m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	// 得到fork出来进程映射在外部namespace的pid
	fmt.Printf("%v", cmd.Process.Pid)

	// 在系统默认创建挂载了memory subsystem的hierarchy上创建cgroup
	os.Mkdir(path.Join(cGroupMemoryHierarchyMount, "testmemorylimit"), 0755)

	// 将容器进程加入到这个cgroup中
	ioutil.WriteFile(path.Join(cGroupMemoryHierarchyMount, "testmemorylimit", "tasks"),
		[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)

	// 限制cgroup进程内存使用
	ioutil.WriteFile(path.Join(cGroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"),
		[]byte("100m"), 0644)

	cmd.Process.Wait()
}
