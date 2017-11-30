package playbook

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/golang/glog"
)

func GetComponents(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("get %s's version error: %v", path, err)
	}

	components := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() && strings.HasSuffix(f.Name(), PlaybookSuffix) {
			components = append(components, strings.TrimSuffix(f.Name(), PlaybookSuffix))
		}
	}

	return components, nil
}

func GetVersions(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("get %s's version error: %v", path, err)
	}

	versions := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		versions = append(versions, f.Name())
	}

	return versions, nil
}

func getInherentProperties(dir string) map[string]interface{} {
	bs, err := ioutil.ReadFile(path.Join(dir, "inherent.yaml"))
	if err != nil {
		glog.Warningf("read file dir inherent property error: %v", err)
		return nil
	}

	value := make(map[string]interface{})
	if err := yaml.Unmarshal(bs, &value); err != nil {
		glog.Warningf("unmarshal inherent error: %s", err)
		return nil
	}

	return value
}
