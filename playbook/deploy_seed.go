package playbook

import (
	"fmt"
	"path"

	"gitee.com/wisecloud/wise-deploy/database"
)

type DeploySeed2 map[string]*Component

func NewDeploySeed(c *database.Cluster, workDir string) *DeploySeed2 {
	cs := DeploySeed2(make(map[string]*Component))
	for _, cp := range c.Components {
		hosts := make(map[string][]*database.Host)
		for k, v := range cp.Hosts {
			hs := make([]*database.Host, 0, len(v))
			for _, h := range v {
				host := getEntireHost(c, h)
				hs = append(hs, host)
			}

			hosts[k] = hs
		}

		cs[cp.Name] = &Component{
			MetaComponent: cp.MetaComponent,
			Hosts:         hosts,
		}

		getInherentProperties(
			path.Join(workDir, cp.Name+PlaybookSuffix, cp.Version),
			cs[cp.Name],
		)
	}

	return &cs
}

func (ds *DeploySeed2) AllHosts() map[string]*database.Host {
	hosts := make(map[string]*database.Host)
	for _, v := range map[string]*Component(*ds) {
		for _, hv := range v.Hosts {
			for _, h := range hv {
				hosts[h.IP] = h
			}
		}
	}

	return hosts
}

type Component struct {
	database.MetaComponent
	Inherent map[string]interface{}
	Hosts    map[string][]*database.Host
}

// func setComponentHost(
// 	clusterID string,
// 	sourceCp *database.Component,
// 	destCp *Component,
// 	workDir string,
// ) error {
// 	destCp.MetaComponent = sourceCp.MetaComponent

// 	destCp.Inherent = getInherentProperties(
// 		path.Join(workDir, sourceCp.Name+PlaybookSuffix, sourceCp.Version),
// 	)

// 	return ConvertHosts(clusterID, sourceCp.Hosts, destCp.Hosts)
// }

func ConvertHosts(
	clusterID string,
	sourceHosts map[string][]string,
	destHost map[string][]*database.Host,
) error {
	for k, v := range sourceHosts {
		hosts := make([]*database.Host, 0, len(v))

		for _, h := range v {
			hh, err := database.Default().RetrieveHost(clusterID, h)
			if err != nil {
				return err
			}
			hosts = append(hosts, hh)
		}

		destHost[k] = hosts
	}

	return nil
}

func getEntireHost(c *database.Cluster, host string) *database.Host {
	for _, h := range c.Hosts {
		if host == h.ID {
			return h
		}
	}

	panic(fmt.Sprintf("can't find the host: %s", host))
}
