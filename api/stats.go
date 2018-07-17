package api

import (
	"luckgo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TotalMasterDbConnections(c *gin.Context) {
	count := model.Srv.SqlSupplier.TotalMasterDbConnections()
	c.JSON(http.StatusOK, gin.H{
		"result": count,
		"code":   "200",
	})
}

func TotalReadDbConnections(c *gin.Context) {
	count := model.Srv.SqlSupplier.TotalReadDbConnections()
	c.JSON(http.StatusOK, gin.H{
		"result": count,
		"code":   "200",
	})
}

func TotalSearchDbConnections(c *gin.Context) {
	count := model.Srv.SqlSupplier.TotalSearchDbConnections()
	c.JSON(http.StatusOK, gin.H{
		"result": count,
		"code":   "200",
	})
}
