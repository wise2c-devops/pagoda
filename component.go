package main

import "gitee.com/wisecloud/wise-deploy/database"

var (
	ComponentMap = map[string]int{
		"registry":     0,
		"etcd":         1,
		"mysql":        2,
		"loadbalancer": 3,
		"k8smaster":    4,
		"k8snode":      5,
		"wisecloud":    6,
	}
)

type ByName []string

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	return ComponentMap[s[i]] < ComponentMap[s[j]]
}

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
