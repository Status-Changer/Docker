package subsystems

// ResourceConfig 用于传递资源限制配置，包含内存限制、CPU时间片权重与核心数目
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// Subsystem 每个Subsystem可以实现以下的4个接口
//
// 这里将cgroup抽象成了path，这是由于cgroup在hierarchy下的路径就是虚拟文件系统中的虚拟路径
type Subsystem interface {
	Name() string

	// Set 设置某个cgroup在这个Subsystem中的资源限制
	Set(path string, res *ResourceConfig) error

	// Apply 将进程添加到某个cgroup中
	Apply(path string, pid int) error

	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&MemorySubsystem{},
	}
)
