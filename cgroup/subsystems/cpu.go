package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/tristan/myrealdocker/log"
)

var logger = log.NewLogger()

type CpuSubSystem struct {
}

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	logger.Infof("cpu cgroup path: [%v]", subsysCgroupPath)
	if err == nil {
		if res.CpuShare != "" {
			err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644)
			if err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	logger.Infof("apply cpu for pid:[%v], cgroupPath:[%v]", pid, cgroupPath)
	subsysCgroupPath, err := GetCgroupPath(s.Name(),cgroupPath, false)
	if err == nil {
		err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}
	
func (s *CpuSubSystem) Name() string {
	return "cpu"
}