package playbook

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/golang/glog"
)

func GetComponents(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("get %s's available components error: %v", path, err))
	}

	components := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() && strings.HasSuffix(f.Name(), PlaybookSuffix) {
			components = append(components, strings.TrimSuffix(f.Name(), PlaybookSuffix))
		}
	}

	return components
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

func getInherentProperties(dir string, cp *Component) {
	buf := &bytes.Buffer{}

	bs, err := ioutil.ReadFile(path.Join(dir, "inherent.yaml"))
	if err != nil {
		glog.Warningf("read file dir inherent property error: %v", err)
		return
	}

	tp := template.Must(template.New("inherent").Funcs(fns).Parse(string(bs)))
	err = tp.Execute(buf, cp)
	tp.Option("missingkey=zero")
	if err != nil {
		glog.Warningf(fmt.Sprintf("execute template for %s/inherent.yaml error: %v", dir, err))
		return
	}

	value := make(map[string]interface{})
	if err := yaml.Unmarshal(buf.Bytes(), &value); err != nil {
		glog.Warningf("unmarshal inherent error: %s", err)
	}

	cp.Inherent = value
}

func GetOrderedComponents(dir string) ([]string, error) {
	inFile, err := os.Open(path.Join(dir, "components_order.conf"))
	if err != nil {
		return nil, fmt.Errorf("read components order error: %v", err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	ret := make([]string, 0)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, nil
}
