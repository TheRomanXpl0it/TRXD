package infos

type InstanceInfo struct {
	Host         string
	InternalPort *int32
	ExternalPort *int32
	Envs         string
	MaxMemory    int32
	MaxCpu       string
	NetName      string
	NetID        string
}
