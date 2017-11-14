package main

import (
	"net/http"

	"gitee.com/wisecloud/wise-deploy/database"

	"github.com/gin-gonic/gin"
)

func retrieveLogs(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	logs, err := database.Default().RetrieveLogs(clusterID)
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

	status, err := clusterRuntime.RetrieveStatus(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.IndentedJSON(http.StatusOK, status)
}
