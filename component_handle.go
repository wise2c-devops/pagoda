package main

import (
	"net/http"

	"github.com/wise2c-devops/pagoda/database"

	"github.com/gin-gonic/gin"
)

func retrieveComponents(c *gin.Context) {
	clusterID := c.Param("cluster_id")

	cs, err := database.Instance().RetrieveComponents(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ccs := make([]*Component, 0, len(cs))
	for _, cp := range cs {
		cc := NewComponent(clusterID, cp)
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

	err := database.Instance().CreateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		cc := NewComponent(clusterID, cp)
		c.JSON(http.StatusOK, cc)
	}
}

func deleteComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentID := c.Param("component_id")

	err := database.Instance().DeleteComponent(clusterID, componentID)
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

	if componentID != cp.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "two component id must be equal",
		})
		return
	}

	err := database.Instance().UpdateComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		cc := NewComponent(clusterID, cp)
		c.JSON(http.StatusOK, cc)
	}
}

func retrieveComponent(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	componentID := c.Param("component_id")

	cp, err := database.Instance().RetrieveComponent(clusterID, componentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cc := NewComponent(clusterID, cp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, cc)
	}
}
