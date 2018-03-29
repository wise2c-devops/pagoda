package playbook

import "testing"
import "github.com/wise2c-devops/pagoda/database"

func TestPreparePlaybooks(t *testing.T) {
	ds := &DeploySeed{
		// "registry": &Component{
		// 	Hosts: map[string][]*database.Host{
		// 		"self": []*database.Host{
		// 			&database.Host{IP: "192.168.0.101"},
		// 		},
		// 	},
		// },
		// "etcd": &Component{
		// 	Hosts: map[string][]*database.Host{
		// 		"self": []*database.Host{
		// 			&database.Host{IP: "192.168.0.101"},
		// 		},
		// 	},
		// },
		// "mysql": &Component{
		// 	Hosts: map[string][]*database.Host{
		// 		"self": []*database.Host{
		// 			&database.Host{IP: "192.168.0.101"},
		// 		},
		// 	},
		// },
		// "loadbalancer": &Component{
		// 	Hosts: map[string][]*database.Host{
		// 		"self": []*database.Host{
		// 			&database.Host{IP: "192.168.0.101"},
		// 		},
		// 	},
		// },
		"kubernetes": &Component{
			Hosts: map[string][]*database.Host{
				"master": []*database.Host{
					&database.Host{IP: "192.168.0.101"},
				},
			},
			MetaComponent: database.MetaComponent{
				Version: "v1.8.6",
			},
			Inherent: map[string]interface{}{
				"endpoint": "192.168.10.1",
			},
		},
		// "wisecloud": &Component{
		// 	Hosts: map[string][]*database.Host{
		// 		"self": []*database.Host{
		// 			&database.Host{IP: "192.168.0.101"},
		// 		},
		// 	},
		// },
	}

	if err := PreparePlaybooks(".", ds); err != nil {
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
