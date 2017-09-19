package main

import "gitee.com/wisecloud/wise-deploy/database"

type Component struct {
	ID       string
	Name     string
	Property map[string]interface{}
	Hosts    []*database.Host
}

func NewComponent(clusterID string, cp *database.Component) (*Component, error) {
	c := &Component{
		ID:       cp.ID,
		Name:     cp.Name,
		Property: cp.Property,
		Hosts:    make([]*database.Host, 0, len(cp.Hosts)),
	}

	for _, h := range cp.Hosts {
		hh, err := database.Instance(sqlConfig).RetrieveHost(clusterID, h)
		if err != nil {
			return nil, err
		}

		c.Hosts = append(c.Hosts, hh)
	}

	return c, nil
}
