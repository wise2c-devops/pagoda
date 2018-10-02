package playbook

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
)

// Template src is template path and dest is template output file
type Template struct {
	Src  string `json:"-"`
	Dest string `json:"-"`
}

const (
	//PlaybookSuffix - suffix for playbook folder
	PlaybookSuffix  = "-playbook"
	ansibleGroupDir = "group_vars"
	hostsTmpl       = "hosts.gotmpl"
	hostsFile       = "hosts"
	tmplDir         = "yat"
	tmplSuffix      = ".gotmpl"
)

func PreparePlaybooks(dir string, ds *DeploySeed) error {
	for k, v := range map[string]*Component(*ds) {
		if err := preparePlaybook(path.Join(dir, k+PlaybookSuffix, v.Version), ds); err != nil {
			return err
		}
	}

	return nil
}

func preparePlaybook(name string, ds *DeploySeed) error {
	tps, err := getTemplatePath(name)
	if err != nil {
		return err
	}

	for _, tp := range tps {
		if err = applyTemplate(tp, ds); err != nil {
			return err
		}
	}

	return nil
}

func applyTemplate(t *Template, ds *DeploySeed) error {
	file, err := os.OpenFile(t.Dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("create template dest file %s error: %s", t.Dest, err)
	}
	defer file.Close()

	content, err := ioutil.ReadFile(t.Src)
	if err != nil {
		return fmt.Errorf("read template src file error: %v", err)
	}

	tp := template.Must(template.New("ansible").Funcs(fns).Parse(string(content)))
	// tp.Option("missingkey=zero")
	err = tp.Execute(file, ds)
	if err != nil {
		return fmt.Errorf("execute template for %s error: %v", t.Dest, err)
	}

	return nil
}

// getTemplatePath - check playbook, get every template path & output file
func getTemplatePath(name string) ([]*Template, error) {
	if err := checkPlaybook(name); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path.Join(name, tmplDir))
	if err != nil {
		return nil, fmt.Errorf("read dir %s error: %v", name, err)
	}

	tps := make([]*Template, 0, len(files))
	for _, f := range files {
		if f.Name() == hostsTmpl {
			t := &Template{
				Src:  path.Join(name, tmplDir, f.Name()),
				Dest: path.Join(name, hostsFile),
			}
			tps = append(tps, t)
		} else {
			t := &Template{
				Src:  path.Join(name, tmplDir, f.Name()),
				Dest: path.Join(name, ansibleGroupDir, strings.TrimSuffix(f.Name(), tmplSuffix)),
			}
			tps = append(tps, t)
		}
	}

	return tps, nil
}

func checkPlaybook(name string) error {
	files, err := ioutil.ReadDir(path.Join(name, tmplDir))
	if err != nil {
		return fmt.Errorf("check %s error: %v", name, err)
	}

	var hasGroupVars, hasHostsGotmpl bool
	d, err := os.Stat(path.Join(name, ansibleGroupDir))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if d.IsDir() {
		hasGroupVars = true
	}

	for _, f := range files {
		if f.Name() != hostsTmpl && !hasGroupVars {
			return fmt.Errorf("have group vars template but have not group directory")
		}

		if f.Name() == hostsTmpl {
			hasHostsGotmpl = true
		}
	}

	if !hasHostsGotmpl {
		return fmt.Errorf("have not %s", hostsTmpl)
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

var fns = template.FuncMap{
	"notLast": func(x int, a interface{}) bool {
		return x < reflect.ValueOf(a).Len()-1
	},
}
