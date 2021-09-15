package subsystems

import (
	"os"
	"path"
	"testing"
)

func TestMemoryCgroup(t *testing.T) {
	memSubSys := MemorySubSystem{}
	resouceConfig := ResourceConfig{
		MemoryLimit: "100m",
	}
	testCgroup := "testmemorylimit"

	err := memSubSys.Set(testCgroup, &resouceConfig)
	if err != nil {
		t.Fatalf("cgroup fail %v", err)
	}

	stat, _ := os.Stat(path.Join(FindCgroupMountpoint("memory"), testCgroup))
	t.Logf("cgroup stat: %+v", stat)

	err = memSubSys.Apply(testCgroup, os.Getgid())
	if err != nil {
		t.Fatalf("cgroup apply %v", err)
	}

	// err = memSubSys.Apply("", os.Getgid())
	// if err != nil {
	// 	t.Fatalf("cgroup Apply %v", err)
	// }

	err = memSubSys.Remove(testCgroup); 
	if err != nil {
		t.Fatalf("cgroup remove %v", err)
	}
}