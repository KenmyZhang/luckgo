package api

import (
	"luckgo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(c *gin.Context) {
	cfg := model.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"result": cfg,
		"code":   "200 OK",
	})
}

func GetVersionDetails(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{"CurrentVersion": model.CurrentVersion, "BuildDate": model.BuildDate, "BuildHash": model.BuildHash},
		"code":   "200 OK",
	})
}