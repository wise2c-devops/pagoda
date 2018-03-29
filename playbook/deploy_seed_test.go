package playbook

import (
	"reflect"
	"testing"

	"github.com/wise2c-devops/pagoda/database"
)

func TestNewDeploySeed(t *testing.T) {
	type args struct {
		c       *database.Cluster
		workDir string
	}
	tests := []struct {
		name string
		args args
		want *DeploySeed
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				c: &database.Cluster{
					Components: []*database.Component{
						&database.Component{
							MetaComponent: database.MetaComponent{
								Name:    "kubernetes",
								Version: "v1.8.6",
							},
						},
					},
				},
				workDir: ".",
			},
			want: &DeploySeed{
				"kubernetes": &Component{
					MetaComponent: database.MetaComponent{
						Name: "kubernetes",
					},
					Inherent: map[string]interface{}{
						"endpoint": 123,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeploySeed(tt.args.c, tt.args.workDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeploySeed's kubernetes = %v", (*got)["kubernetes"].Inherent)
				t.Errorf("NewDeploySeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
