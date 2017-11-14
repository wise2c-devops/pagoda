package main

import "gitee.com/wisecloud/wise-deploy/database"
import "reflect"

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

	hmap := make(map[string]*database.Host)
	for _, h := range cp.Hosts {
		hh, err := database.Default().RetrieveHost(clusterID, h)
		if err != nil {
			return nil, err
		}

		c.Hosts = append(c.Hosts, hh)
		hmap[h] = hh
	}

	for k, v := range cp.Property {
		if t := reflect.TypeOf(v); t.Kind() != reflect.Slice {
			continue
		}

		hosts := v.([]interface{})

		for _, ih := range hosts {
			h := ih.(string)
			c.Property[k] = hmap[h]
		}
	}

	return c, nil
}
