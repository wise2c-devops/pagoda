package database

import (
	"testing"

	"gitee.com/wisecloud/wise-deploy/cluster"
)

func TestCluster(t *testing.T) {
	config := &EngineConfig{
		SQLType: "sqlite3",
		ShowSQL: true,
	}

	e, err := NewEngine(config)
	if err != nil {
		t.Error(err)
	}

	c1 := &cluster.Cluster{
		ID:          "c1",
		Name:        "c1",
		Description: "c1",
	}
	err = e.CreateCluster(c1)
	if err != nil {
		t.Error(err)
	}

	c2 := &cluster.Cluster{
		ID:          "c2",
		Name:        "c2",
		Description: "c2",
	}
	err = e.CreateCluster(c2)

	c3 := &cluster.Cluster{
		ID:          "c3",
		Name:        "c3",
		Description: "c3",
	}
	err = e.CreateCluster(c3)

	cs, err := e.RetrieveClusters()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cs)
	}
}
