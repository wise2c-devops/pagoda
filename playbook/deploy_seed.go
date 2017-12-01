package playbook

import (
	"fmt"
	"path"

	"gitee.com/wisecloud/wise-deploy/database"
)

type DeploySeed2 map[string]*Component

func (ds *DeploySeed2) AllHosts() []*database.Host {
	hosts := make([]*database.Host, 0)
	for _, v := range map[string]*Component(*ds) {
		for _, hv := range v.Hosts {
			hosts = append(hosts, hv...)
		}
	}

	return hosts
}

// func (d *DeploySeed2) EsEndpoint() string {
// 	if esVip, find := d.LoadBalancer.Property["es_vip"]; find {
// 		return esVip.(string)
// 	}

// 	if len(d.K8sNode.Hosts) > 0 {
// 		return d.K8sNode.Hosts[0].IP
// 	}

// 	return ""
// }

// func (d *DeploySeed2) OtherEndpoint() string {
// 	if otherVip, find := d.LoadBalancer.Property["other_vip"]; find {
// 		return otherVip.(string)
// 	}

// 	if len(d.K8sNode.Hosts) > 0 {
// 		return d.K8sNode.Hosts[0].IP
// 	}

// 	return "192.168.0.1"
// }

// func (d *DeploySeed2) RegistryEndpoint() string {
// 	if registryVip, find := d.LoadBalancer.Property["registry_vip"]; find {
// 		return registryVip.(string)
// 	}

// 	if len(d.Registry.Hosts) > 0 {
// 		return d.Registry.Hosts[0].IP
// 	}

// 	return ""
// }

// func (d *DeploySeed2) MySQLEndpoint() string {
// 	if mysqlVip, find := d.LoadBalancer.Property["mysql_vip"]; find {
// 		return mysqlVip.(string)
// 	}

// 	if len(d.MySQL.Hosts) > 0 {
// 		return d.MySQL.Hosts[0].IP
// 	}

// 	return ""
// }

// func (d *DeploySeed2) K8sEndpoint() string {
// 	if k8sVip, find := d.LoadBalancer.Property["k8s_vip"]; find {
// 		return k8sVip.(string)
// 	}

// 	if len(d.K8sMaster.Hosts) > 0 {
// 		return d.K8sMaster.Hosts[0].IP
// 	}

// 	return ""
// }

type Component struct {
	database.MetaComponent
	Inherent map[string]interface{}
	Hosts    map[string][]*database.Host
}

func NewDeploySeed(c *database.Cluster, workDir string, components []string) *DeploySeed2 {
	ds := DeploySeed2(make(map[string]*Component))

	for _, cp := range components {
		component := getComponent(c, cp, workDir)
		ds[cp] = component
	}

	return &ds
}

func setComponentHost(
	clusterID string,
	sourceCp *database.Component,
	destCp *Component,
	workDir string,
) error {
	destCp.MetaComponent = sourceCp.MetaComponent

	destCp.Inherent = getInherentProperties(
		path.Join(workDir, sourceCp.Name+PlaybookSuffix, sourceCp.Version),
	)

	return ConvertHosts(clusterID, sourceCp.Hosts, destCp.Hosts)
}

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

func getComponent(c *database.Cluster, componentName string, workDir string) *Component {
	for _, cp := range c.Components {
		if componentName == cp.Name {
			hosts := make(map[string][]*database.Host)
			for k, v := range cp.Hosts {
				hs := make([]*database.Host, 0, len(v))
				for _, h := range v {
					host := getEntireHost(c, h)
					hs = append(hs, host)
				}

				hosts[k] = hs
			}
			return &Component{
				MetaComponent: cp.MetaComponent,
				Hosts:         hosts,
				Inherent:      getInherentProperties(path.Join(workDir, cp.Name+PlaybookSuffix, cp.Version)),
			}
		}
	}

	panic(fmt.Sprintf("can't find the component: %s", componentName))
}

func getEntireHost(c *database.Cluster, host string) *database.Host {
	for _, h := range c.Hosts {
		if host == h.ID {
			return h
		}
	}

	panic(fmt.Sprintf("can't find the host: %s", host))
}
