package utils

const (
	DefaultContainerIdLength            int    = 32
	ContainerStatusRunning              string = "running"
	ContainerStatusStopped              string = "stopped"
	ContainerStatusExited               string = "exited"
	DefaultContainerInfoStorageLocation string = "/var/run/Docker/%s/"
	DefaultContainerInfoConfigName      string = "config.json"
	DefaultContainerLogFileName         string = "container.log"
)
