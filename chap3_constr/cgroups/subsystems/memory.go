package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
}

// Set 设置cgroupPath对应的cgroup内存限制
func (s *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			// 设置cgroup的内存限制，也就是写入到对应目录的memory.limit_in_bytes文件中
			if err := ioutil.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"),
				[]byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgruop memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

// Remove 删除cgroupPath对应的cgroup
func (s *MemorySubsystem) Remove(cgroupPath string) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsystemCgroupPath)
	} else {
		return err
	}
}

// Apply 把一个进程加入到cgroupPath对应的cgroup中
func (s *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	if SubsystemsCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(SubsystemsCgroupPath, "tasks"),
			[]byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (s *MemorySubsystem) Name() string {
	return "memory"
}
