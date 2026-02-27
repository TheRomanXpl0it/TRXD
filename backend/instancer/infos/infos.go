package infos

type InstanceInfo struct {
	Name         string
	Domain       string
	UseDomain    bool
	InternalPort *int32
	ExternalPort *int32
	Envs         string
	MaxMemory    int32
	MaxCpu       string
	NetID        string
	Labels       map[string]string
}
