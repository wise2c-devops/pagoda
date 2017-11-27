package playbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func getVersion(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("get %s's version error: %v", path, err)
	}

	versions := make([]string, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if strings.HasPrefix(f.Name(), "v") {
			versions = append(versions, f.Name())
		}
	}

	return versions, nil
}

func getInherentProperties(dir string) (map[string]interface{}, error) {
	bs, err := ioutil.ReadFile(path.Join(dir, "inherent.yaml"))
	if err != nil {
		return nil, err
	}

	value := make(map[string]interface{})
	if err := json.Unmarshal(bs, &value); err != nil {
		return nil, err
	}

	return value, nil
}
