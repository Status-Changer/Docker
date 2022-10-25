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
	mountURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkspace(rootURL, mountURL)
	cmd.Dir = mountURL
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

func NewWorkspace(rootURL, mountUrl string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mountUrl)
}

// CreateReadOnlyLayer 把busybox.tar解压到busybox目录下，作为只读层
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exists, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
		return
	}
	if !exists {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
			return
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL,
			"-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error. %v", busyboxTarURL, err)
		}
	}
}

// CreateWriteLayer 创建writeLayer文件夹作为容器唯一的可写层
func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}

// CreateMountPoint 创建mnt文件夹作为挂载点，并完成只读层和可写层的挂载
func CreateMountPoint(rootURL, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
		return
	}

	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
