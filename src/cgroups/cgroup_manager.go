package cgroups

import (
	"Docker/src/cgroups/subsystems"
	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	// Path cgroup在hierarchy中的路径
	Path string

	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{Path: path}
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subsystemIns := range subsystems.Instance {
		subsystemIns.Set(c.Path, res)
	}
	return nil
}

// Apply 将进程PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subsystemIns := range subsystems.Instance {
		subsystemIns.Apply(c.Path, pid)
	}
	return nil
}

// Destroy 释放各个subsystem挂载中的cgroup
func (c *CgroupManager) Destroy() error {
	for _, subsystemIns := range subsystems.Instance {
		if err := subsystemIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
