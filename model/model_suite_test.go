package model_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/onsi/gomega"
	"luckgo/log"
	"luckgo/model"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Model Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	//创建server实例
	srv := model.NewServer()
	// 加载配置文件
	model.LoadConfig("config.json")

	// 初始化 logger
	srv.Log = log.NewLogger(model.Cfg.LoggerConfigFromLoggerConfig())

	// 将golang中默认的 logger重定向到这个指定的server logger
	log.RedirectStdLog(srv.Log)

	// 使用server logger 作为全局的logger
	log.InitGlobalLogger(srv.Log)

	// 初始化数据库
	srv.SqlSupplier = model.NewSqlSupplier()
})

var _ = AfterSuite(func() {
    model.Srv.ShutDown()
})
