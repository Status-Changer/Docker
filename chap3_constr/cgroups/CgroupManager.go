package cgroups

import "Docker/chap3_constr/cgroups/subsystems"

type CgroupManager struct {
	// Path cgroup在hierarchy中的路径
	Path string

	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{Path: path}
}

// Apply 将进程PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subsystemIns := range subsystems.SubsystemsIns {
		subsystemIns.Apply(c.Path, pid)
	}
	return nil
}
