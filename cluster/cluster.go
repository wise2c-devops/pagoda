package cluster

type Cluster struct {
	ID          string `xorm:"varchar(255) notnull pk 'id'"`
	Name        string `xorm:"varchar(255) notnull 'name'"`
	Description string `xorm:"varchar(255) 'description'"`
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
