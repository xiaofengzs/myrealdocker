package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

/* */
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logger.Errorf("Exec loop path error %v", err)
		return err
	}
	logger.Infof("Find path %s", path)

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logger.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logger.Errorf("init read pipe error %v", err)
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("Get current location err %v", err)
		return 
	}
	logger.Infof("Current location is %s", pwd)

	pivotRoot(pwd)

	defaultMpountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMpountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	err = os.Mkdir(pivotDir, 0777)
	if err != nil {
		return err
	}

	err = syscall.PivotRoot(root, pivotDir)
	if err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	err = syscall.Chdir("/")
	if err != nil {
		return fmt.Errorf("chdir / err %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	err = syscall.Unmount(pivotDir, syscall.CLONE_DETACHED)
	if err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)
}
