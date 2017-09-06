package cluster

type Cluster struct {
	ID          string
	Name        string
	Description string
	Hosts       []*Host
	Components  []*Component
}

type Host struct {
	ID          string
	HostName    string
	IP          string
	Description string
}

type Component struct {
	Name     string
	Property map[string]interface{}
	Hosts    []string
}
