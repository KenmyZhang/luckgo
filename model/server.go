package model

import (
	"luckgo/log"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GracefulServer struct {
	Server      *http.Server
	Router      *gin.Engine
	SqlSupplier *SqlSupplier
	Log         *log.Logger
}

var Srv *GracefulServer

func NewServer() *GracefulServer {
	log.Info("new server")
	Srv = &GracefulServer{}
	return Srv
}

func (gracefulServer *GracefulServer) Start() error {
	log.Info("server start")
	gracefulServer.Server = &http.Server{
		Addr:         Cfg.ServiceSettings.ListenAddress,
		Handler:      gracefulServer.Router,
		ReadTimeout:  time.Duration(Cfg.ServiceSettings.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(Cfg.ServiceSettings.WriteTimeout) * time.Second,
	}
	log.Info("Start Listening and serving HTTP on " + Cfg.ServiceSettings.ListenAddress)
	log.Info(fmt.Sprintf("Current version is %v (%v/%v)", CurrentVersion, BuildDate, BuildHash))
	if err := gracehttp.Serve(gracefulServer.Server); err != nil {
		return err
	}
	return nil
}

func (gracefulServer *GracefulServer) ShutDown() {
	if gracefulServer.Server != nil {
		gracefulServer.Server.Close()
	}
	if gracefulServer.SqlSupplier != nil {
		gracefulServer.SqlSupplier.Close()
	}
	log.Info("shutdown server successfully")
}
