package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	if err == nil {
		if res.MemoryLimit != "" {
			err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
			if err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			} 
		}
		return nil
	}  else {
		return err
	}
}

func (s *MemorySubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err ==  nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroup, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err == nil {
		err := ioutil.WriteFile(path.Join(subsysCgroup, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}

