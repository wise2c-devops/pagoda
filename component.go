package main

import "gitee.com/wisecloud/wise-deploy/database"
import "gitee.com/wisecloud/wise-deploy/playbook"

type Component struct {
	database.MetaComponent
	Hosts map[string][]*database.Host `json:"hosts"`
}

func NewComponent(clusterID string, cp *database.Component) *Component {
	c := &Component{
		MetaComponent: cp.MetaComponent,
		Hosts:         playbook.ConvertHosts(clusterID, cp.Hosts),
	}

	return c
}
