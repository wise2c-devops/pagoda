package cluster

type Cluster struct {
	ID          string       `xorm:"varchar(255) notnull pk 'id'"`
	Name        string       `xorm:"varchar(255) notnull 'name'"`
	Description string       `xorm:"varchar(255) 'description'"`
	Hosts       []*Host      `xorm:"-"`
	Component   []*Component `xorm:"-"`
}

type ClusterHost struct {
	ClusterID string `xorm:"varchar(255) notnull 'id'"`
	Host      *Host  `xorm:"json notnull"`
}

type Host struct {
	ID          string
	HostName    string
	IP          string
	Description string
}

type ClusterComponent struct {
	ClusterID string     `xorm:"varchar(255) notnull 'id'"`
	Component *Component `xorm:"json notnull"`
}

type Component struct {
	Name     string
	Property map[string]interface{}
	Hosts    []string
}
