package main

import (
	"fmt"
	"log"
	"net/http"
	"onlineShopping/models"
	"onlineShopping/pkg/setting"
	"onlineShopping/routers"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	gin.SetMode(setting.TotalConfig.Server.RunMode)
	endPoint := fmt.Sprintf("%s:%s", setting.TotalConfig.Server.Host, setting.TotalConfig.Server.Port)

	server := &http.Server{
		Addr:         endPoint,
		Handler:      routers.InitRouter(),
		ReadTimeout:  setting.TotalConfig.Server.ReadTimeout,
		WriteTimeout: setting.TotalConfig.Server.WriteTimeout,
		// ReadTimeout:    10000000,
		// WriteTimeout:   10000000,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[info] start http server listening %s", endPoint)
	server.ListenAndServe()
}
