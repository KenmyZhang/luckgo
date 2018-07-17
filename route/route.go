package route

import (
	"luckgo/model"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Root           *gin.Engine
	Helloworld *gin.RouterGroup
	Stats          *gin.RouterGroup
	Prometheus     *gin.RouterGroup
}

var BaseRouter *Router

func InitRoute(router *gin.Engine) {
	BaseRouter := &Router{Root: router}
	model.Srv.Router = router
	BaseRouter.HelloWorld()
	BaseRouter.InitStats()
	BaseRouter.InitPrometheus()
	BaseRouter.InitConfig()
	return
}
