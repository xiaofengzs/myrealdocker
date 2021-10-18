package main

import (
	"os"
	"strings"

	"github.com/tristan/myrealdocker/container"
)

/*
 */
 func Run(tty bool, comArray []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logger.Errorf("New parent process error")
		return 
	}
	if err:= parent.Start(); err != nil {
		logger.Error(err)
	}

	sendInitCommand(comArray, writePipe)
	parent.Wait()

	mntURL := "/home/tristan/mnt/"
	rootURL := "/home/tristan/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logger.Infof("command all is [%s]", command)
	writePipe.WriteString(command)
	writePipe.Close()
}