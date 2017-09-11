package main

import (
	"net/http"

	"gitee.com/wisecloud/wise-deploy/cluster"
	"gitee.com/wisecloud/wise-deploy/database"

	"github.com/gin-gonic/gin"
)

func retrieveComponents(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	cs, err := database.Instance(sqlConfig).RetrieveComponents(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cs)
	}
}

func createComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	cp := &cluster.Component{}
	if err := c.BindJSON(cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err := database.Instance(sqlConfig).CreateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

func deleteComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentName := c.Param("component_name")

	err := database.Instance(sqlConfig).DeleteComponent(clusterID, componentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

func updateComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	cp := &cluster.Component{}
	if err := c.BindJSON(cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err := database.Instance(sqlConfig).UpdateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

func retrieveComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentName := c.Param("component_name")

	cp, err := database.Instance(sqlConfig).RetrieveComponent(clusterID, componentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cp)
	}
}
