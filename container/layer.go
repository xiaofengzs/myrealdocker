package container

import (
	"fmt"
	"os"
	"os/exec"
)

func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	// busyboxURL := rootURL + "busybox/"
	busyboxURL := fmt.Sprintf("%sbusybox/", rootURL)
	busyboxTarURL := fmt.Sprintf("%sbusybox.tar/", rootURL)

	exist, err := PathExists(busyboxURL)
	if err != nil {
		logger.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if !exist {
		err := os.Mkdir(busyboxURL, 0777)
		if err != nil {
			logger.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
		}
		_, err = exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput()
		if err != nil {
			logger.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := fmt.Sprintf("%swriteLayer/", rootURL)
	logger.Infof("WriteLayer path is %v", writeURL)
	err := os.Mkdir(writeURL, 0777)
	if err != nil {
		logger.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	err := os.Mkdir(mntURL, 0777)
	if err != nil {
		logger.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}
	logger.Infof("rootURL is [%v]", rootURL)
	logger.Infof("mntURL is [%v]", mntURL)
	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + "busybox"
	logger.Infof("dirs path is %v", dirs)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	logger.Infof("mount command [%v]", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		logger.Errorf("Mount writeable layer and init layer, err: %v", err)
	}
}

func DeleteWorkSpace(rootURL string, mntURL string) {
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.Errorf("Unmount AUFS fs failed, err: %v", err)
	}
	err = os.RemoveAll(mntURL)
	if err != nil {
		logger.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLater/"
	err := os.RemoveAll(writeURL)
	if err != nil {
		logger.Errorf("Remove dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
