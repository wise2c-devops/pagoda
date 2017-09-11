package database

import (
	"fmt"

	"gitee.com/wisecloud/wise-deploy/cluster"

	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	//just init
	_ "github.com/mattn/go-sqlite3"
)

type SQLEngine struct {
	xe *xorm.Engine
}

type EngineConfig struct {
	SQLType      string
	ShowSQL      bool
	ShowExecTime bool
}

var (
	i *SQLEngine
)

func Instance(config *EngineConfig) *SQLEngine {
	if i == nil {
		var err error
		i, err = NewEngine(config)
		if err != nil {
			panic(fmt.Sprintf("get sql engine instance error %v", err))
		}

		glog.V(3).Info("create sql engine instance")
	}

	return i
}

func NewEngine(config *EngineConfig) (*SQLEngine, error) {
	e := &SQLEngine{}
	if config.SQLType == "sqlite3" {
		engine, err := xorm.NewEngine("sqlite3", "cluster.db")
		if err != nil {
			return nil, err
		}
		e.xe = engine
	}

	e.xe.ShowSQL(config.ShowSQL)
	e.xe.ShowExecTime(config.ShowExecTime)
	if err := e.xe.Sync2(new(cluster.Cluster)); err != nil {
		return nil, err
	}

	if err := e.xe.Sync2(new(cluster.ClusterComponent)); err != nil {
		return nil, err
	}

	if err := e.xe.Sync2(new(cluster.ClusterHost)); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *SQLEngine) RetrieveClusters() ([]*cluster.Cluster, error) {
	clusters := make([]*cluster.Cluster, 0)
	err := e.xe.Find(&clusters)
	return clusters, err
}

func (e *SQLEngine) CreateCluster(c *cluster.Cluster) error {
	c.ID = newUUID()
	c.State = cluster.Initial

	if _, err := e.xe.InsertOne(c); err != nil {
		return err
	}

	return nil
}

func (e *SQLEngine) DeleteCluster(id string) error {
	_, err := e.xe.Exec("delete from cluster where id = ?", id)
	return err
}

func (e *SQLEngine) UpdateCluster(c *cluster.Cluster) error {
	i, err := e.xe.Id(c.ID).Update(c)
	if err != nil {
		return err
	}

	if i == 0 {
		return fmt.Errorf("can't find cluster %s", c.ID)
	}

	return nil
}

func (e *SQLEngine) RetrieveCluster(id string) (c *cluster.Cluster, err error) {
	c = &cluster.Cluster{}
	_, err = e.xe.ID(id).Get(c)
	if err != nil {
		glog.V(3).Info(err)
		return
	}

	cs, err := e.RetrieveComponents(id)
	if err != nil {
		glog.V(3).Info(err)
		return
	}
	c.Components = cs

	hs, err := e.RetrieveHosts(id)
	if err != nil {
		glog.V(3).Info(err)
		return
	}
	c.Hosts = hs

	return
}

func (e *SQLEngine) RetrieveComponents(
	clusterID string,
) (cs []*cluster.Component, err error) {
	ccs := make([]*cluster.ClusterComponent, 0)
	if err = e.xe.Find(&ccs); err != nil {
		return
	}

	cs = make([]*cluster.Component, 0, len(ccs))
	for _, cc := range ccs {
		cs = append(cs, cc.Component)
	}
	return
}

func (e *SQLEngine) CreateComponent(clusterID string, cp *cluster.Component) error {
	cc := &cluster.ClusterComponent{
		ClusterID:     clusterID,
		ComponentName: cp.Name,
		Component:     cp,
	}

	_, err := e.xe.InsertOne(cc)
	return err
}

func (e *SQLEngine) DeleteComponent(clusterID string, name string) error {
	_, err := e.xe.Exec(
		"delete from cluster_component where cluster_id = ? and component_name = ?",
		clusterID,
		name,
	)
	return err
}

func (e *SQLEngine) UpdateComponent(clusterID string, cp *cluster.Component) error {
	ccp := &cluster.ClusterComponent{
		ClusterID:     clusterID,
		ComponentName: cp.Name,
		Component:     cp,
	}
	_, err := e.xe.Update(ccp)
	return err
}

func (e *SQLEngine) RetrieveComponent(
	clusterID string,
	name string,
) (*cluster.Component, error) {
	cc := &cluster.ClusterComponent{
		ClusterID:     clusterID,
		ComponentName: name,
	}
	_, err := e.xe.Get(cc)
	return cc.Component, err
}

func (e *SQLEngine) RetrieveHosts(
	clusterID string,
) (hs []*cluster.Host, err error) {
	chs := make([]*cluster.ClusterHost, 0)
	if err = e.xe.Find(&chs); err != nil {
		return
	}

	hs = make([]*cluster.Host, 0, len(hs))
	for _, ch := range chs {
		hs = append(hs, ch.Host)
	}
	return
}

func (e *SQLEngine) CreateHost(clusterID string, h *cluster.Host) error {
	h.ID = newUUID()
	ch := &cluster.ClusterHost{
		ClusterID: clusterID,
		HostID:    h.ID,
		Host:      h,
	}

	_, err := e.xe.InsertOne(ch)
	return err
}

func (e *SQLEngine) DeleteHost(clusterID string, hostID string) error {
	_, err := e.xe.Exec(
		"delete from cluster_host where cluster_id = ? and host_id = ?",
		clusterID,
		hostID,
	)
	return err
}

func (e *SQLEngine) UpdateHost(clusterID string, h *cluster.Host) error {
	ch := &cluster.ClusterHost{
		ClusterID: clusterID,
		HostID:    h.ID,
		Host:      h,
	}
	_, err := e.xe.Update(ch)
	return err
}

func (e *SQLEngine) RetrieveHost(
	clusterID string,
	hostID string,
) (*cluster.Host, error) {
	ch := &cluster.ClusterHost{
		ClusterID: clusterID,
		HostID:    hostID,
	}
	_, err := e.xe.Get(ch)
	return ch.Host, err
}
