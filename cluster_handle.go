package main

import (
	"net/http"

	"gitee.com/wisecloud/wise-deploy/database"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func retrieveClusters(c *gin.Context) {
	cs, err := database.Default().RetrieveClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cs)
	}
}

func createCluster(c *gin.Context) {
	cluster := &database.Cluster{}
	if err := c.BindJSON(cluster); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := database.Default().CreateCluster(cluster)
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

	if err := database.Default().DeleteCluster(clusterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

func updateCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	cluster := &database.Cluster{}
	if err := c.BindJSON(cluster); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if clusterID != cluster.ID {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "two cluster id must be equal",
		})
		return
	}

	err := database.Default().UpdateCluster(cluster)
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
		Default().
		RetrieveCluster(clusterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}

func ErrResponse(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors[len(c.Errors)-1]
	render.WriteJSON(
		c.Writer,
		gin.H{
			"error": err.Error(),
		},
	)
}
