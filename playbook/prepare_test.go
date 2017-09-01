package playbook

import "testing"

func TestPreparePlaybooks(t *testing.T) {
	ds := &DeploySeed{
		Hosts: Hosts{
			Etcd: []string{
				"172.20.8.1",
				"172.20.8.2",
				"172.20.8.3",
			},
			K8sMaster: []string{
				"172.20.8.4",
				"172.20.8.5",
				"172.20.8.6",
			},
			K8sNode: K8sNode{
				Wise2cController: []string{
					"172.20.8.7",
					"172.20.8.8",
					"172.20.8.9",
				},
				Normal: []string{
					"172.20.8.10",
					"172.20.8.11",
					"172.20.8.12",
				},
			},
			LoadBalancer: []string{
				"172.20.8.13",
				"172.20.8.14",
				"172.20.8.15",
			},
			Registry: []string{
				"172.20.8.16",
				"172.20.8.17",
				"172.20.8.18",
			},
			MysqlMaster: "172.20.8.19",
			MysqlSlave1: "172.20.8.20",
			MysqlSlave2: "172.20.8.21",
		},
		Vips: Vips{
			Interface: "etch0",
			NetMask:   16,
			K8s:       "172.20.9.1",
			Es:        "172.20.9.2",
			Other:     "172.20.9.3",
			Registry:  "172.20.9.4",
		},
	}

	if err := PreparePlaybooks("/home/mian/workspace/k8s/", ds); err != nil {
		t.Error(err)
	}
}

func TestGetTemplatePath(t *testing.T) {
	tps, err := getTemplatePath("good-playbook")
	if err != nil {
		t.Error(err)
	}

	for _, tp := range tps {
		t.Log(tp.Src)
		t.Log(tp.Dest)
	}

	_, err = getTemplatePath("bad0-playbook")
	if err == nil {
		t.Error("bad-playbook should have error")
	} else {
		t.Log(err)
	}

	_, err = getTemplatePath("bad1-playbook")
	if err == nil {
		t.Error("bad-playbook should have error")
	} else {
		t.Log(err)
	}

	_, err = getTemplatePath("bad2-playbook")
	if err == nil {
		t.Error("bad-playbook should have error")
	} else {
		t.Log(err)
	}
}
