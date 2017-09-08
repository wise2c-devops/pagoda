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
		for _, c := range cs {
			t.Log(c)
		}
	}
}

func TestComponent(t *testing.T) {
	config := &EngineConfig{
		SQLType: "sqlite3",
		ShowSQL: true,
	}

	e, err := NewEngine(config)
	if err != nil {
		t.Error(err)
	}

	c := &cluster.Component{
		Name: "etcd",
		Property: map[string]interface{}{
			"caFile": "ca.crt",
		},
		Hosts: []string{
			"aaa",
			"bbb",
		},
	}
	err = e.CreateComponent("f4a27554-41c6-4a6b-bd30-e13c131756c1", c)
	if err != nil {
		t.Error(err)
	}

	c.Property["caKey"] = "ca.key"
	err = e.UpdateComponent("f4a27554-41c6-4a6b-bd30-e13c131756c1", c)
	if err != nil {
		t.Error(err)
	}

	c1, err := e.RetrieveComponent("f4a27554-41c6-4a6b-bd30-e13c131756c1", "etcd")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(c1)
	}
}
