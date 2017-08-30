package playbook

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Template struct {
	Src  string `json:"-"`
	Dest string `json:"-"`
}

func getTemplatePath(name string) ([]*Template, error) {
	if err := checkPlaybook(name); err != nil {
		return nil, err
	}

	return nil, nil
}

func checkPlaybook(name string) error {
	files, err := ioutil.ReadDir(name)
	if err != nil {
		return fmt.Errorf("read dir %s error: %v", name, err)
	}

	var hasGroupVars, hasHostsGotmpl bool
	d, err := os.Stat(path.Join(name, "group_vars"))
	if err != nil {
		return err
	} else if d.IsDir() {
		hasGroupVars = true
	}

	for _, f := range files {
		if f.Name() != "hosts.gotmpl" && !hasGroupVars {
			return fmt.Errorf("have group vars template but have not group_vars")
		}

		if f.Name() == "hosts.gotmpl" {
			hasHostsGotmpl = true
		}
	}
	if !hasHostsGotmpl {
		return fmt.Errorf("have not hosts.gotmpl")
	}

	return nil
}

func GetFileFromDir(dir string, cf func(os.FileInfo) bool) (fs []os.FileInfo, err error) {
	files, err := ioutil.ReadDir(dir)

	for _, f := range files {
		if cf(f) {
			fs = append(fs, f)
		}
	}

	return
}
