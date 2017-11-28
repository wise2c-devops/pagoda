package main

import "gitee.com/wisecloud/wise-deploy/database"
import "gitee.com/wisecloud/wise-deploy/playbook"

type Component struct {
	database.MetaComponent
	Hosts map[string][]*database.Host `json:"hosts"`
}

func NewComponent(clusterID string, cp *database.Component) (*Component, error) {
	c := &Component{
		MetaComponent: cp.MetaComponent,
		Hosts:         make(map[string][]*database.Host),
	}

	if err := playbook.ConvertHosts(clusterID, cp.Hosts, c.Hosts); err != nil {
		return nil, err
	}

	return c, nil
}
