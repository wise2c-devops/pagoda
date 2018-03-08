package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path"

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
	ansibleChan = make(chan *database.Notification, 5)

	workDir *string
)

func init() {
	workDir = flag.String("w", ".", "ansible playbook should be placed in it")
	database.DBPath = flag.String("db-path", "/deploy/cluster.db", "sqlite db file path")
	database.InitSQL = flag.String("init-sql", "/root/table.sql", "sql file to init db")
	flag.Parse()

	runtime.Run(*workDir)
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
		v1.GET("/components", components)
		v1.GET("/components/:component_name/versions", versions)
		v1.GET("/components/:component_name/properties/:version", properties)

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

func versions(c *gin.Context) {
	component := c.Param("component_name")
	versions, err := playbook.GetVersions(
		path.Join(
			*workDir,
			component+playbook.PlaybookSuffix,
		),
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, versions)
	}
}

func properties(c *gin.Context) {
	component := c.Param("component_name")
	version := c.Param("version")
	c.Header("Content-Type", "application/json; charset=utf-8")

	http.ServeFile(
		c.Writer,
		c.Request,
		path.Join(
			*workDir,
			fmt.Sprintf(
				"%s-playbook/%s/properties.json",
				component,
				version,
			),
		),
	)
}

func install(c *gin.Context) {
	op := &runtime.LaunchParameters{}
	if err := c.BindJSON(&op); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clusterID := c.Param("cluster_id")

	cluster, err := database.Instance().RetrieveCluster(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := runtime.StartOperate(cluster, op); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cluster.State = database.Processing
	err = database.Instance().UpdateCluster(cluster)
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
	runtime.StopOperate()
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
	runtime.Notify(config)

	c.Status(http.StatusOK)
}

func stats(c *gin.Context) {
	wc, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		glog.Errorf("upgrade error: %v", err)
		return
	}
	defer wc.Close()

	ch := runtime.Register(c.Request.RemoteAddr)
	defer runtime.Annul(c.Request.RemoteAddr)

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
