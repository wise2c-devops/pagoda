package main

import "gitee.com/wisecloud/wise-deploy/database"

type Component struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Property map[string]interface{} `json:"properties"`
	Hosts    []*database.Host       `json:"hosts"`
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

type Notification struct {
	Data map[string]interface{} `json:"data"`
	Now  string                 `json:"now"`
	Task struct {
		Name  string `json:"name"`
		State string `json:"state"`
	} `json:"task"`
	Stage string `json:"stage"`
	State string `json:"state"`
	Host  string `json:"host"`
}
