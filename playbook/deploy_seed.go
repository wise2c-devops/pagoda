package playbook

import (
	"path"

	"gitee.com/wisecloud/wise-deploy/database"
)

type DeploySeed2 struct {
	Registry     *Component
	Etcd         *Component
	MySQL        *Component
	LoadBalancer *LoadBalancer
	K8sMaster    *Component
	K8sNode      *Component
	WiseCloud    *Component
	Hosts        []*database.Host
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

func NewDeploySeed(c *database.Cluster, workDir string) *DeploySeed2 {
	ds := &DeploySeed2{
		Registry:     &Component{},
		Etcd:         &Component{},
		MySQL:        &Component{},
		LoadBalancer: &LoadBalancer{},
		K8sMaster:    &Component{},
		K8sNode:      &Component{},
		WiseCloud:    &Component{},
	}
	ds.Hosts = c.Hosts

	for _, cp := range c.Components {
		switch cp.Name {
		case "etcd":
			setComponentHost(c.ID, cp, ds.Etcd, workDir)
		case "registry":
			setComponentHost(c.ID, cp, ds.Registry, workDir)
		case "mysql":
			setComponentHost(c.ID, cp, ds.MySQL, workDir)
		case "loadbalancer":
			setComponentHost(c.ID, cp, (*Component)(ds.LoadBalancer), workDir)
		case "k8smaster":
			setComponentHost(c.ID, cp, ds.K8sMaster, workDir)
		case "k8snode":
			setComponentHost(c.ID, cp, ds.K8sNode, workDir)
		case "wisecloud":
			setComponentHost(c.ID, cp, ds.WiseCloud, workDir)
		}
	}

	return ds
}

func setComponentHost(
	clusterID string,
	sourceCp *database.Component,
	destCp *Component,
	workDir string,
) error {
	destCp.MetaComponent = sourceCp.MetaComponent

	inherent, err := getInherentProperties(
		path.Join(workDir, sourceCp.Name+PlaybookSuffix, sourceCp.Version),
	)
	if err != nil {
		return err
	}
	destCp.Inherent = inherent

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
