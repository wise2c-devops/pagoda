package main

import (
	"net/http"

	"gitee.com/wisecloud/wise-deploy/cluster"
	"gitee.com/wisecloud/wise-deploy/database"

	"github.com/gin-gonic/gin"
)

func retrieveClusters(c *gin.Context) {
	cs, err := database.Instance(sqlConfig).RetrieveClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cs)
	}
}

func createCluster(c *gin.Context) {
	cluster := &cluster.Cluster{}
	if err := c.BindJSON(cluster); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := database.Instance(sqlConfig).CreateCluster(cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}

func deleteCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	if err := database.Instance(sqlConfig).DeleteCluster(clusterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

func updateCluster(c *gin.Context) {
	cluster := &cluster.Cluster{}
	if err := c.BindJSON(cluster); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := database.Instance(sqlConfig).UpdateCluster(cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}

func retrieveCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	if cluster, err := database.
		Instance(sqlConfig).
		RetrieveCluster(clusterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}
