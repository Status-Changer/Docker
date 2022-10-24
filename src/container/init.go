package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// RunContainerInitProcess 在容器内部执行，使用mount去挂载proc文件系统
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	setUpMount()

	// exec.LookPath在系统的PATH环境变量中寻找命令的绝对路径
	// 因此可以直接通过命令名字，而不是路径，如ls而不是/bin/ls来调用命令
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)

	// 容器创建之后，执行的第一个进程是init进程而不是用户进程，这和预期是不一样的
	// 但是因为init进程的PID=1，如果kill掉则容器就退出了，因此在这里使用syscall.Exec
	// syscall.Exec最终调用了execve这个系统调用，执行当前filename对应的程序，并覆盖掉当前进程的镜像、数据及堆栈信息
	// 因此运行ps -ef时会发现/bin/bash是PID=1的进程
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

// readUserCommand 读取管道中传入的消息
func readUserCommand() []string {
	// 下标3是读管道，详见container_process.go中的cmd.ExtraFiles
	pipe := os.NewFile(uintptr(3), "pipe")

	// 如果父进程还没有写入内容，此时管道会等待输入
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

// pivotRoot 类似于pivot_root系统调用，改变当前的root文件系统，使得新root和原先的root不在同一个文件系统下
func pivotRoot(root string) error {
	// bind mount是把相同的内容换了一个挂载点
	if err := syscall.Mount(root, root, "bind",
		syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	// 此时old_root挂载在rootfs/.pivot_root上，切换到新的rootfs
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// 修改当前的目录到'/'
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// pivot_root会移除之前对root文件系统的依赖，因此可以unmount之前的文件系统
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// 删除临时的文件夹
	return os.Remove(pivotDir)
}

func setUpMount() {
	// 获得当前路径
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	// MS_NOEXEC是指本文件系统中不允许运行其他程序
	// MS_NOSUID是指本系统运行时不允许set userID或set groupID
	// MS_NODEV从Linux 2.4以来所有mount系统默认指定的参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs",
		syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
