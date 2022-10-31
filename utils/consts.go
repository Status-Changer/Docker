package utils

const (
	DefaultContainerIdLength            int    = 32
	ContainerStatusRunning              string = "running"
	ContainerStatusStopped              string = "stopped"
	ContainerStatusExited               string = "exited"
	DefaultContainerInfoStorageLocation string = "/var/run/Docker/%s/"
	DefaultContainerInfoConfigName      string = "config.json"
	EnvironmentExecPid                  string = "docker_pid"
	EnvironmentExecCommand              string = "docker_cmd"
)
