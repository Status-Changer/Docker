package main

import (
	"Docker/chap3_constr/cgroups/subsystems"
	"Docker/chap3_constr/container"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in the container. Do not call it outside",
	// 获取传递过来的command参数；执行容器初始化操作
	Action: func(context *cli.Context) error {
		log.Infof("init called")
		err := container.RunContainerInitProcess()
		return err
	},
}

// flags是运行时使用--指定的参数
var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit Docker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	// 判断参数是否含有command；获取用户的command；调用Run准备启动容器
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpuset"),
			CpuSet:      context.String("cpushare"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}
