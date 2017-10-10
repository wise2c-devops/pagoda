package database

import (
	"time"
)

const (
	Initial    = "initial"
	Processing = "processing"
	Success    = "success"
	Failed     = "failed"
)

type Cluster struct {
	ID          string       `xorm:"varchar(255) notnull pk 'id'" json:"id"`
	Name        string       `xorm:"varchar(255) notnull unique 'name'" json:"name"`
	Description string       `xorm:"varchar(255) 'description'" json:"description"`
	State       string       `xorm:"varchar(255) 'state'" json:"state"`
	Hosts       []*Host      `xorm:"-" json:"hosts,omitempty"`
	Components  []*Component `xorm:"-" json:"components,omitempty"`
}

type ClusterHost struct {
	ClusterID string `xorm:"varchar(255) notnull pk 'cluster_id'"`
	HostID    string `xorm:"varchar(255) notnull pk 'host_id'"`
	IP        string `xorm:"varchar(25) notnull unique 'ip'"`
	Hostname  string `xorm:"varchar(255) notnull unique 'hostname'"`
	Host      *Host  `xorm:"json notnull 'host'"`
}

type Host struct {
	ID          string `json:"id"`
	HostName    string `json:"hostname"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type ClusterComponent struct {
	ClusterID     string     `xorm:"varchar(255) notnull pk 'cluster_id'"`
	ComponentID   string     `xorm:"varchar(255) notnull pk 'component_id'"`
	ComponentName string     `xorm:"varchar(255) notnull 'component_name'"`
	Component     *Component `xorm:"json notnull 'component'"`
}

type Component struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Property map[string]interface{} `json:"properties"`
	Hosts    []string               `json:"hosts"`
}

type ClusterLog struct {
	ClusterID string        `xorm:"varchar(255) notnull 'cluster_id'"`
	Created   time.Time     `xorm:"created"`
	Log       *Notification `xorm:"json notnull 'log'"`
}

type Notification struct {
	Data map[string]interface{} `json:"data"`
	Now  string                 `json:"time"`
	Task struct {
		Name  string `json:"name"`
		State string `json:"state"`
	} `json:"task"`
	Stage string `json:"stage"`
	State string `json:"state"`
	Host  string `json:"host"`
}
