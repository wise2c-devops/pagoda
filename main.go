package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"

	"gitee.com/wisecloud/wise-deploy/database"
	"gitee.com/wisecloud/wise-deploy/playbook"

	"gopkg.in/yaml.v2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ansibleChan = make(chan map[string]interface{})
	commands    = NewCommands()

	workDir = flag.String("w", ".", "ansible playbook should be placed in it")

	sqlConfig = &database.EngineConfig{
		SQLType:      "sqlite3",
		ShowSQL:      true,
		ShowExecTime: true,
	}
)

func init() {
	flag.Parse()
	go commands.Launch(*workDir)
}

func main() {
	defer glog.Flush()

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowMethods = append(config.AllowMethods, "DELETE")
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.StaticFile("/favicon.ico", "favicon.ico")

	v1 := r.Group("/v1")
	{
		v1.GET("/clusters", retrieveClusters)
		v1.POST("/clusters", createCluster)
		v1.DELETE("/clusters/:cluster_id", deleteCluster)
		v1.PUT("/clusters/:cluster_id", updateCluster)
		v1.GET("/clusters/:cluster_id", retrieveCluster)

		v1.GET("/clusters/:cluster_id/hosts", retrieveHosts)
		v1.POST("/clusters/:cluster_id/hosts", createHost)
		v1.DELETE("/clusters/:cluster_id/hosts/:host_id", deleteHost)
		v1.PUT("/clusters/:cluster_id/hosts/:host_id", updateHost)
		v1.GET("/clusters/:cluster_id/hosts/:host_id", retrieveHost)

		v1.GET("/clusters/:cluster_id/components", retrieveComponents)
		v1.POST("/clusters/:cluster_id/components", createComponent)
		v1.DELETE("/clusters/:cluster_id/components/:component_id", deleteComponent)
		v1.PUT("/clusters/:cluster_id/components/:component_id", updateComponent)
		v1.GET("/clusters/:cluster_id/components/:component_id", retrieveComponent)

		v1.PUT("/config", setConfig)
		v1.GET("/config", getConfig)
		v1.PUT("/install/:cluster_id", install)
		v1.PUT("/reset", reset)
		v1.PUT("/stop", stop)
		v1.POST("/notify", notify)
		v1.GET("/stats", stats)
		v1.Static("/docs", "apidoc")
	}

	// Listen and Server in 0.0.0.0:8080
	r.Run("0.0.0.0:8080")
}

func setConfig(c *gin.Context) {
	config := &playbook.DeploySeed{}
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

func install(c *gin.Context) {
	op := make(map[string]string)
	if err := c.BindJSON(&op); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clusterID := c.Param("cluster_id")

	cluster, err := database.
		Instance(sqlConfig).
		RetrieveCluster(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	config := playbook.NewDeploySeed(cluster)

	if err := playbook.PreparePlaybooks(*workDir, config); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	commands.Install(step)

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func reset(c *gin.Context) {
	// config := &playbook.DeploySeed{}
	// if err := c.BindJSON(config); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// if err := playbook.PreparePlaybooks(*workDir, config); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// commands.Reset(config.Step)

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func stop(c *gin.Context) {
	commands.Stop()
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

func saveConfig(s *playbook.DeploySeed) error {
	// b, err := yaml.Marshal(s)
	// if err != nil {
	// 	return err
	// }

	// err = ioutil.WriteFile("init.yaml", b, os.ModePerm)
	// if err != nil {
	// 	return err
	// }

	// err = playbook.PreparePlaybooks(*workDir, s)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func readConfig(filename string) (*playbook.DeploySeed, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorf("read init.yml error: %v", err)
		return nil, err
	}

	d := &playbook.DeploySeed{}
	err = yaml.Unmarshal(b, d)
	if err != nil {
		glog.Errorf("unmarshal init.yml error: %v", err)
		return nil, err
	}

	return d, nil
}
