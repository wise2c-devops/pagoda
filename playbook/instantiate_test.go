package playbook

import (
	"testing"

	"gitee.com/wisecloud/wise-deploy/database"
)

func TestInstantiate(t *testing.T) {
	c := &database.Cluster{
		Name: "test",
		ID:   "1",
		Components: []*database.Component{
			&database.Component{
				MetaComponent: database.MetaComponent{
					Name: "etcd",
				},
			},
			&database.Component{
				MetaComponent: database.MetaComponent{
					Name: "mysql",
				},
			},
		},
	}

	if err := InstantiateCluster("/home/mian/workspace/wise2c-playbook/", c); err != nil {
		t.Error(err.Error())
	} else {
		t.Log("hehe")
	}
}
