package cgroup

import (
	"github.com/tristan/myrealdocker/cgroup/subsystems"
	"github.com/tristan/myrealdocker/log"
)

var logger = log.NewLogger()

type CgroupMamager struct {
	Path string
	Resource *subsystems.CpuSubSystem
}

func NewCgroupManager(path string) *CgroupMamager {
	return &CgroupMamager{}
}

func (c *CgroupMamager) Apply(pid int) error {
	for _, subSysIns := range(subsystems.SubsystemIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupMamager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range(subsystems.SubsystemIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupMamager) Destory() error {
	for _, subSysIns := range(subsystems.SubsystemIns) {
		err := subSysIns.Remove(c.Path)
		if err != nil {
			logger.Warnf("remove cgrouo fail %v", err)
		}
	}
	return nil
}