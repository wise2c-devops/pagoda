package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"gitee.com/wisecloud/wise-deploy/database"
	"gitee.com/wisecloud/wise-deploy/playbook"
	"gitee.com/wisecloud/wise-deploy/runtime"

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
	ansibleChan    = make(chan *database.Notification, 5)
	commands       = runtime.NewCommands()
	clusterRuntime = runtime.NewClusterRuntime()

	workDir = flag.String("w", ".", "ansible playbook should be placed in it")
)

func init() {
	flag.Parse()
	go commands.Launch(*workDir)
	go clusterRuntime.Run()
}

func main() {
	defer glog.Flush()

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowMethods = append(config.AllowMethods, "DELETE")
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.StaticFile("/favicon.ico", "/root/favicon.ico")

	v1 := r.Group("/v1")

	{
		for k := range runtime.ComponentMap {
			r.Group("/v1").StaticFile(
				fmt.Sprintf(
					"/components/%s/properties",
					k,
				),
				fmt.Sprintf(
					"/%s/%s-playbook/file/properties.json",
					*workDir,
					k,
				),
			)
		}

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

		v1.GET("/clusters/:cluster_id/logs", retrieveLogs)
		v1.GET("/clusters/:cluster_id/status", retrieveClusterStatus)

		v1.PUT("/clusters/:cluster_id/deployment", install)
		v1.DELETE("/clusters/:cluster_id/deployment", stop)
		v1.POST("/notify", notify)
		v1.GET("/stats", stats)

		v1.Static("/docs", "apidoc")
	}

	// Listen and Server in 0.0.0.0:8080
	r.Run("0.0.0.0:8080")
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

	cluster, err := database.Default().RetrieveCluster(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	config := playbook.NewDeploySeed(cluster)

	if err = playbook.PreparePlaybooks(*workDir, config); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	o, ok := op["operation"]
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "please give me a operation",
		})
		return
	}
	if o == "install" {
		commands.Install(cluster)
	} else if o == "reset" {
		commands.Reset(cluster)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "please give me a right operation",
		})
		return
	}

	if err = commands.Process(); err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	cluster.State = database.Processing
	err = database.Default().UpdateCluster(cluster)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "started",
	})
}

func stop(c *gin.Context) {
	commands.Stop()
}

func notify(c *gin.Context) {
	config := &database.Notification{}
	if err := c.BindJSON(config); err != nil {
		glog.Error(err)
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

	ch := clusterRuntime.Registe(c.Request.RemoteAddr)
	defer clusterRuntime.Unregiste(c.Request.RemoteAddr)

	for {
		m := <-ch
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
