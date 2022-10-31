package main

import (
	"Docker/cgroups/subsystems"
	"Docker/container"
	"Docker/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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

		createTty := context.Bool("ti")
		detach := context.Bool("d")
		if createTty && detach {
			return fmt.Errorf("ti and d parameter cannot be both provided")
		}

		volume := context.String("v")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpuset"),
			CpuSet:      context.String("cpushare"),
		}
		log.Infof("createTty=%v", createTty)
		containerName := context.String("name")
		Run(createTty, cmdArray, resConf, volume, containerName)
		return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		imageName := context.Args().Get(0)
		commitContainer(imageName)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all containers",
	Action: func(context *cli.Context) error {
		ListAllContainers()
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "execute a command into container",
	Action: func(context *cli.Context) error {
		// for callback 如果已经指定了环境变量，说明C代码已经运行过了，直接返回就行
		if os.Getenv(utils.EnvironmentExecPid) != "" {
			log.Infof("pid callback pid %d", os.Getpid())
			return nil
		}

		// 期望的调用方式是./docker exec 容器名 command
		if len(context.Args()) < 2 {
			return fmt.Errorf("missing container name or command")
		}

		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}
