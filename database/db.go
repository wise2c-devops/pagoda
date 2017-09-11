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

func (e *SQLEngine) RetrieveClusters() (clusters []*cluster.Cluster, err error) {
	err = e.xe.Find(&clusters)
	return
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
	_, err := e.xe.Update(c)
	return err
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
	c.Components = make([]*cluster.Component, 0, len(cs))
	for _, cc := range cs {
		c.Components = append(c.Components, cc.Component)
	}

	hs, err := e.RetrieveHosts(id)
	if err != nil {
		glog.V(3).Info(err)
		return
	}
	c.Hosts = make([]*cluster.Host, 0, len(hs))
	for _, ch := range hs {
		c.Hosts = append(c.Hosts, ch.Host)
	}

	return
}

func (e *SQLEngine) RetrieveComponents(
	clusterID string,
) (ccs []*cluster.ClusterComponent, err error) {
	err = e.xe.Find(&ccs)
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
) (chs []*cluster.ClusterHost, err error) {
	err = e.xe.Find(&chs)
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
