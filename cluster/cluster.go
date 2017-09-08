package cluster

const (
	Initial = "initial"
	Ongoing = "ongoing"
	Success = "success"
	Failed  = "failed"
)

type Cluster struct {
	ID          string       `xorm:"varchar(255) notnull pk 'id'"`
	Name        string       `xorm:"varchar(255) notnull 'name'"`
	Description string       `xorm:"varchar(255) 'description'"`
	State       string       `xorm:"varchar(255) 'state'"`
	Hosts       []*Host      `xorm:"-"`
	Components  []*Component `xorm:"-"`
}

type ClusterHost struct {
	ClusterID string `xorm:"varchar(255) notnull 'id'"`
	Host      *Host  `xorm:"json notnull 'host'"`
}

type Host struct {
	ID          string
	HostName    string
	IP          string
	Description string
}

type ClusterComponent struct {
	ClusterID string     `xorm:"varchar(255) notnull 'id'"`
	Component *Component `xorm:"json notnull 'component'"`
}

type Component struct {
	Name     string
	Property map[string]interface{}
	Hosts    []string
}
