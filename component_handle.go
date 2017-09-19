package main

import (
	"net/http"

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
		return
	}

	ccs := make([]*Component, 0, len(cs))
	for _, cp := range cs {
		cc, err := NewComponent(clusterID, cp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ccs = append(ccs, cc)
	}

	c.JSON(http.StatusOK, ccs)
}

func createComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	cp := &database.Component{}
	if err := c.BindJSON(cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := database.Instance(sqlConfig).CreateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cp)
	}
}

func deleteComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentID := c.Param("component_id")

	err := database.Instance(sqlConfig).DeleteComponent(clusterID, componentID)
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
	componentID := c.Param("component_id")

	cp := &database.Component{}
	if err := c.BindJSON(cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if componentID != cp.Name {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "two component name must be equal",
		})
		return
	}

	err := database.Instance(sqlConfig).UpdateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cp)
	}
}

func retrieveComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentID := c.Param("component_id")

	cp, err := database.Instance(sqlConfig).RetrieveComponent(clusterID, componentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cc, err := NewComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cc)
	}
}
