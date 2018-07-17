package main

import (
	"luckgo/log"
	"luckgo/model"
	"luckgo/route"
	"luckgo/tools"
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var config = flag.String("config", "config.json", "path to the config file")

func main() {
	flag.Parse()

	//创建server实例
	srv := model.NewServer()
	defer srv.ShutDown()

	// 加载配置文件
	model.LoadConfig(*config)

	// 初始化 logger
	srv.Log = log.NewLogger(model.Cfg.LoggerConfigFromLoggerConfig())

	// 将golang中默认的 logger重定向到这个指定的server logger
	log.RedirectStdLog(srv.Log)

	// 使用server logger 作为全局的logger
	log.InitGlobalLogger(srv.Log)

	//设置运行模式
	router := gin.New()
	if model.Cfg.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	//注册日志记录器
	router.Use(tools.MiddleLogger(tools.DefaultMetricPath))
	//注册 Prometheus 监控
	router.Use(tools.Prometheus())
	// 初始化注册路由
	route.InitRoute(router)

	// 初始化数据库
	srv.SqlSupplier = model.NewSqlSupplier()

	//启动服务
	err := srv.Start()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("Server on " + model.Cfg.ServiceSettings.ListenAddress + " stopped")

	os.Exit(0)
}
