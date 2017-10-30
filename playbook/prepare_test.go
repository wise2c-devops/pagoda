package playbook

import "testing"
import "gitee.com/wisecloud/wise-deploy/database"

func TestPreparePlaybooks(t *testing.T) {
	ds := &DeploySeed2{
		Registry: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		Etcd: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		MySQL: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		LoadBalancer: &LoadBalancer{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		K8sMaster: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		K8sNode: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
		WiseCloud: &Component{
			Hosts: []*database.Host{
				&database.Host{IP: "192.168.0.101"},
			},
		},
	}

	if err := PreparePlaybooks("/home/mian/workspace/wise2c-playbook/", ds); err != nil {
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
