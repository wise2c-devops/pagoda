package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

type Template struct {
	Src  string `json:"-"`
	Dest string `json:"-"`
}

type DeploySeed struct {
	Hosts     Hosts
	Templates []*Template `yaml:"-" json:"-"`
	Vips      Vips
}

func main() {
	flag.Parse()
	defer glog.Flush()

	r := gin.Default()

	r.StaticFile("/favicon.ico", "./favicon.ico")

	v1 := r.Group("/v1")
	{
		v1.PUT("/config", setConfig)
		v1.GET("/config", getConfig)
		v1.PUT("/launch", launch)
		v1.POST("/notify", notify)
	}

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setConfig(c *gin.Context) {
	config := &DeploySeed{}
	if err := c.BindWith(config, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := saveConfig(config); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func getConfig(c *gin.Context) {
	d, err := readConfig("init.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.IndentedJSON(http.StatusOK, d)
}

func launch(c *gin.Context) {
	cmd := exec.Command("./install.sh", "install")
	err := cmd.Start()
	if err != nil {
		c.Error(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func notify(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func saveConfig(s *DeploySeed) error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("init.yaml", b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(filename string) (*DeploySeed, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorf("read init.yml error: %v", err)
		return nil, err
	}

	d := &DeploySeed{}
	err = yaml.Unmarshal(b, d)
	if err != nil {
		glog.Errorf("unmarshal init.yml error: %v", err)
		return nil, err
	}

	return d, nil
}
