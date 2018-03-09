package main

import (
	"net/http"

	"gitee.com/wisecloud/wise-deploy/database"
	"gitee.com/wisecloud/wise-deploy/playbook"
	"gitee.com/wisecloud/wise-deploy/runtime"

	"github.com/gin-gonic/gin"
)

func retrieveLogs(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	logs, err := database.Instance().RetrieveLogs(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, logs)
	}
}

func retrieveClusterStatus(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	status, err := runtime.RetrieveStatus(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, status)
}

func components(c *gin.Context) {
	ret, err := playbook.GetOrderedComponents()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, ret)
}
