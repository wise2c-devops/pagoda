package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
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

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ansibleChan = make(chan map[string]interface{}, 15)

	workDir = flag.String("w", ".", "ansible playbook should be placed in it")
)

func init() {
	flag.Parse()

	if err := os.Chdir(*workDir); err != nil {
		glog.Fatalf("change working directory to %s error: %v", *workDir, err)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		glog.V(3).Infof("file %s have mode %s", f.Name(), f.Mode().String())
		if f.IsDir() && strings.HasSuffix(f.Name(), "-playbook") {

		}
	}
}

func main() {
	defer glog.Flush()

	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")

	v1 := r.Group("/v1")
	{
		v1.PUT("/config", setConfig)
		v1.GET("/config", getConfig)
		v1.PUT("/launch", launch)
		v1.POST("/notify", notify)
		v1.GET("/stats", stats)
	}

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setConfig(c *gin.Context) {
	config := &DeploySeed{}
	if err := c.BindJSON(config); err != nil {
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
	config := make(map[string]interface{})
	if err := c.BindJSON(&config); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	glog.V(4).Info(config)
	select {
	case ansibleChan <- config:
	default:
	}

	c.Status(http.StatusOK)
}

func stats(c *gin.Context) {
	wc, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		glog.Errorf("upgrade error: %v", err)
		return
	}
	defer wc.Close()

	for {
		m := <-ansibleChan
		b, err := json.Marshal(m)
		if err != nil {
			glog.Errorf("marshal ansible message error: %v", err)
			break
		}

		err = wc.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			glog.Errorf("write message error: %v", err)
			break
		}
	}
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
