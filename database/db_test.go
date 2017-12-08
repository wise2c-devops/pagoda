package database

import (
	"testing"
)

var (
	tSQLConfig = &EngineConfig{
		SQLType:      "sqlite3",
		ShowSQL:      false,
		ShowExecTime: true,
		DBURL:        "./cluster.db",
		InitSQL:      "./table.sql",
	}
	ti = Instance(sqlConfig)

	tc = []struct {
		c          *Cluster
		available  bool
		upState    string
		finalState string
	}{
		{
			&Cluster{
				ID:          newUUID(),
				Name:        "Test",
				State:       Initial,
				Description: "for mian unit test",
			},
			true,
			Success,
			Success,
		},
		{
			&Cluster{
				ID:          newUUID(),
				Name:        "test",
				State:       Initial,
				Description: "for another unit test",
			},
			false,
			Success,
			Success,
		},
	}
)

func TestCluster(t *testing.T) {
	var num int
	cs, err := ti.RetrieveClusters()
	if err != nil {
		t.Error(err)
	}

	if len(cs) != num {
		t.Errorf("except get %d cluster, bug get %d", num, len(cs))
	}

	for _, c := range tc {
		err := ti.CreateCluster(c.c)
		if err != nil {
			t.Error(err)
		}

		if c.available {
			num++
		}
	}

	for _, c := range tc {
		c.c.State = c.upState
		err := ti.UpdateCluster(c.c)
		if (err != nil) != c.available {
			t.Errorf("%v", err)
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

	c := &Component{
		MetaComponent: MetaComponent{
			Name: "etcd",
			Property: map[string]interface{}{
				"key": "value",
			},
		},
		Hosts: map[string][]string{
			"aaa": []string{"aaa"},
			"bbb": []string{"bbb"},
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

func TestHost(t *testing.T) {
	config := &EngineConfig{
		SQLType: "sqlite3",
		ShowSQL: true,
	}

	e, err := NewEngine(config)
	if err != nil {
		t.Error(err)
	}

	h := &Host{
		HostName:    "k8s01",
		IP:          "172.20.8.1",
		Description: "001",
	}

	err = e.CreateHost("f4a27554-41c6-4a6b-bd30-e13c131756c1", h)
	if err != nil {
		t.Error(err)
	}

	h.Description = "002"
	err = e.UpdateHost("f4a27554-41c6-4a6b-bd30-e13c131756c1", h)
	if err != nil {
		t.Error(err)
	}

	h1, err := e.RetrieveHost("f4a27554-41c6-4a6b-bd30-e13c131756c1", h.ID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(h1)
	}
}

func TestNil(t *testing.T) {
	var s *SQLEngine
	if s == nil {
		t.Log("yes")
	} else {
		t.Log("no")
	}
}
