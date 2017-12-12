package database

import (
	"os"
	"testing"
)

func newClusterComponent(name, pName string) *ClusterComponent {
	return &ClusterComponent{
		ComponentName: name,
		Component: &Component{
			MetaComponent: MetaComponent{
				Name:    name,
				Version: "1.12.6",
				Property: map[string]interface{}{
					"key": pName,
				},
			},
		},
	}
}

var (
	tSQLConfig = &EngineConfig{
		SQLType:      "sqlite3",
		ShowSQL:      false,
		ShowExecTime: true,
		DBURL:        "./cluster.db",
		InitSQL:      "./table.sql",
	}

	tc = []struct {
		c         *Cluster
		available bool
		upState   string
	}{
		{
			&Cluster{
				Name:        "Test",
				State:       Initial,
				Description: "for mian unit test",
			},
			true,
			Success,
		},
		{
			&Cluster{
				Name:        "test",
				State:       Initial,
				Description: "for another unit test",
			},
			false,
			Success,
		},
	}

	tcpa = []struct {
		cp        *ClusterComponent
		available bool
		upVersion string
	}{
		{
			newClusterComponent("docker", "docker"),
			true,
			"17.12",
		},
		{
			newClusterComponent("docker", "Docker"),
			false,
			"17.12",
		},
	}

	tcpb = []struct {
		cp        *ClusterComponent
		available bool
		upVersion string
	}{
		{
			newClusterComponent("docker", "Docker"),
			true,
			"17.12",
		},
		{
			newClusterComponent("docker", "docker"),
			false,
			"17.12",
		},
	}

	tcpm = map[string][]struct {
		cp        *ClusterComponent
		available bool
		upVersion string
	}{
		"Test": tcpa,
		"test": tcpb,
	}
)

func TestCluster(t *testing.T) {
	os.Remove(tSQLConfig.DBURL)
	ti := Instance(tSQLConfig)

	var num int
	for _, c := range tc {
		err := ti.CreateCluster(c.c)
		if (err != nil) == c.available {
			t.Errorf("cluster %s create unexpected", c.c.Name)
			return
		}

		if c.available {
			num++
		}

		c.c.State = c.upState
		err = ti.UpdateCluster(c.c)
		if (err != nil) == c.available {
			t.Errorf("cluster %s update unexpected", c.c.Name)
			return
		}

		cc, err := ti.RetrieveCluster(c.c.ID)
		if (err != nil) == c.available {
			t.Errorf("cluster %s retrieve unexpected", c.c.Name)
			return
		}

		if (cc.State != c.upState) == c.available {
			t.Errorf("%s state is identical", cc.Name)
			return
		}
	}

	cs, err := ti.RetrieveClusters()
	if err != nil {
		t.Error(err)
		return
	}

	if len(cs) != num {
		t.Errorf("except get %d cluster, bug get %d", num, len(cs))
		return
	}

	if !t.Run("test component", testComponent) {
		return
	}

	for _, c := range tc {
		err := ti.DeleteCluster(c.c.ID)
		if (err != nil) == c.available {
			t.Errorf("cluster %s delete unexpected", c.c.Name)
			return
		}
	}
}

func testComponent(t *testing.T) {
	ti := Instance(tSQLConfig)

	for _, c := range tc {
		var num int
		for i := 0; i < len(tcpm[c.c.Name]); i++ {
			cp := tcpm[c.c.Name][i]
			err := ti.CreateComponent(c.c.ID, cp.cp.Component)
			if (err != nil) == (c.available && cp.available) {
				t.Errorf("component %s create unexcepted", cp.cp.Component.Property["key"])
				return
			}
			cp.cp.ComponentID = cp.cp.Component.ID

			if c.available && cp.available {
				num++
			}

			cp.cp.Component.Version = cp.upVersion
			err = ti.UpdateComponent(c.c.ID, cp.cp.Component)
			if (err != nil) == (c.available && cp.available) {
				t.Errorf("component %s update unexcepted", cp.cp.Component.Property["key"])
				return
			}

			cpp, err := ti.RetrieveComponent(c.c.ID, cp.cp.ComponentID)
			if (err != nil) == (c.available && cp.available) {
				t.Errorf("component %s retrieve unexcepted", cp.cp.Component.Property["key"])
				return
			}

			if (cpp.Version != cp.upVersion) == (c.available && cp.available) {
				t.Errorf("%s version is identical", cp.cp.Component.Property["key"])
				return
			}
		}

		cs, err := ti.RetrieveComponents(c.c.ID)
		if err != nil {
			t.Error(err)
			return
		}

		if len(cs) != num {
			t.Errorf("except get %d components, bug get %d", num, len(cs))
			return
		}
	}

	for _, c := range tc {
		for _, cp := range tcpm[c.c.Name] {
			err := ti.DeleteComponent(c.c.ID, cp.cp.ComponentID)
			if (err != nil) == (c.available && cp.available) {
				t.Errorf("component %s delete unexpected", cp.cp.Component.Property["key"])
				return
			}
		}
	}
}

func testHost(t *testing.T) {
	config := &EngineConfig{
		SQLType: "sqlite3",
		ShowSQL: false,
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
