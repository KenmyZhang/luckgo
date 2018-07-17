package api

import (
	"luckgo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func Helloworld(c *gin.Context) {
	params := &Params{}
	if c.BindJSON(params) != nil {
		c.JSON(http.StatusOK, model.NewInvalidParamError(model.InvalidParam, "api.Helloworld", "Params", "invalid json body"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
		"code":   200,
		"msg": fmt.Sprintf("hello world, status, %v", params.Status),
	})
}
