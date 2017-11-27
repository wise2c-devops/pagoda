package playbook

import (
	"fmt"
	"io/ioutil"
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
