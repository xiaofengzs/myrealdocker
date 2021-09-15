package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tristan/myrealdocker/app"
	"github.com/tristan/myrealdocker/cgroup"
	"github.com/tristan/myrealdocker/cgroup/subsystems"
	"github.com/tristan/myrealdocker/log"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime inplementation.
The purpose of this project is to learn how docker works and how to write a docker by ourselvers
Enjoy it, just for fun.`

var logger = log.NewLogger()

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the defailt ACSII formatter
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("Run failed %v", err.Error)
	}
}

//  这里定义了runCommand的Flags，其作用类似于运行命令时使用--来指定参数
var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroup limit mydocker run -ti [commad]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpushare"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}

// 这里，定义了initCommand的具体操作，此操作为内部方法，禁止外部调用
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container, DO not call it outside",

	/*
		1. 获取传递过来的command参数
		2. 执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		logger.Info("init come ont")
		cmd := context.Args().Get(0)
		logger.Infof("command in init %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

/*
 */
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logger.Errorf("New parent process error")
		return 
	}
	if err:= parent.Start(); err != nil {
		logger.Error(err)
	}

	cgroupManager := cgroup.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Set(res)
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logger.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

