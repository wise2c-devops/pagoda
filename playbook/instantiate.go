package playbook

import (
	"fmt"
	"os"
	"path"

	"github.com/golang/glog"

	"gitee.com/wisecloud/wise-deploy/database"
)

const (
	clusterTemplate = "cluster-template"
)

var (
	ansibleFile = []string{
		"file",
		"template",
		"scripts",
		"yat",
		"ansible.cfg",
		"clean.ansible",
		"install.ansible",
	}
)

func InstantiateCluster(wd string, cluster *database.Cluster) error {
	oldwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}
	defer os.Chdir(oldwd)

	newFolder := fmt.Sprintf("cluster-%s-%s", cluster.Name, cluster.ID)

	if err := os.Mkdir(
		newFolder,
		0755,
	); err != nil {
		return err
	}

	for _, com := range cluster.Components {
		if err := os.Mkdir(
			path.Join(
				newFolder,
				com.Name+PlaybookSuffix,
			),
			0755,
		); err != nil {
			return err
		}

		for _, f := range ansibleFile {
			p := path.Join(
				"cluster-template",
				com.Name+PlaybookSuffix,
				f,
			)

			if _, err := os.Stat(p); os.IsNotExist(err) {
				glog.Errorf("%s's %s folder is not exist", com.Name, f)
				continue
			} else {
				glog.Errorf("check %s's %s folder error: %v", com.Name, f, err)
			}

			if err := os.Symlink(
				p,
				path.Join(
					newFolder,
					com.Name+PlaybookSuffix,
					f,
				),
			); err != nil {
				return err
			}
		}
	}

	return nil
}
