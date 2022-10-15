package main

import (
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
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
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
	},
	// 判断参数是否含有command；获取用户的command；调用Run准备启动容器
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}
