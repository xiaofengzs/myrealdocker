package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/tristan/myrealdocker/log"
)

var logger = log.NewLogger()

/*
	这里是父进程，也就是当前进程执行的内容，根据上一章介绍的内容，应该比较容易明报。
*/
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logger.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/home/tristan/mnt/"
	rootURL := "/home/tristan/"
	NewWorkSpace(rootURL, mntURL)
	cmd.Dir = mntURL	
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}