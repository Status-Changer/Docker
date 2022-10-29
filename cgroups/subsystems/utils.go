package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// FindCgroupMountPoint 找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// 以memory行为例：
	// 35 25 0:31 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:17 - cgroup cgroup rw,memory
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")

		// 这个例子中是rw,memory
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				// 这个例子中应该返回/sys/fs/cgroup/memory
				return fields[4]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// GetCgroupPath 得到cgroup在文件系统中的绝对路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)
	// 路径已经存在；或者不存在但是创建开关打开
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
