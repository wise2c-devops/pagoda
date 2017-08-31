package playbook

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/golang/glog"
)

type K8sNode struct {
	Wise2cController []string `yaml:"wise2cController"`
	Normal           []string `yaml:"normal"`
}

type Hosts struct {
	Etcd         []string `yaml:"etcd"`
	K8sMaster    []string `yaml:"k8sMaster"`
	K8sNode      K8sNode  `yaml:"k8sNode"`
	LoadBalancer []string `yaml:"loadBalancer"`
	Registry     []string `yaml:"registry"`
	MysqlMaster  string   `yaml:"mysqlMaster"`
	MysqlSlave1  string   `yaml:"mysqlSlave1"`
	MysqlSlave2  string   `yaml:"mysqlSlave2"`
	Distincts    []string `yaml:"-" json:"-"`
}

type Vips struct {
	Interface string
	NetMask   int `yaml:"netMask"`
	K8s       string
	Es        string
	Registry  string `yaml:"registry"`
	Other     string
}

type DeploySeed struct {
	Hosts Hosts
	// Templates []*Template `yaml:"-" json:"-"`
	Vips Vips
}

// Template src is template path and dest is template output file
type Template struct {
	Src  string `json:"-"`
	Dest string `json:"-"`
}

const (
	playbookSuffix  = "-playbook"
	ansibleGroupDir = "group_vars"
	hostsTmpl       = "hosts.gotmpl"
	hostsFile       = "hosts"
	tmplDir         = "yat"
	tmplSuffix      = ".gotmpl"
)

func PreparePlaybooks(dir string, ds *DeploySeed) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		glog.V(3).Infof("file %s have mode %s", f.Name(), f.Mode().String())
		if f.IsDir() && strings.HasSuffix(f.Name(), playbookSuffix) {
			if err = preparePlaybook(path.Join(dir, f.Name()), ds); err != nil {
				return err
			}
		}
	}

	return nil
}

func preparePlaybook(name string, ds *DeploySeed) error {
	tps, err := GetTemplatePath(name)
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
	err = tp.Execute(file, ds)
	if err != nil {
		return fmt.Errorf("execute template for %s error: %v", t.Dest, err)
	}

	return nil
}

// GetTemplatePath - check playbook, get every template path & output file
func GetTemplatePath(name string) ([]*Template, error) {
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
				Src:  path.Join(name, f.Name()),
				Dest: path.Join(name, hostsFile),
			}
			tps = append(tps, t)
		} else {
			t := &Template{
				Src:  path.Join(name, f.Name()),
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
		return err
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
