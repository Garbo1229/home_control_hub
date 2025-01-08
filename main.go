package main

import (
	"home_control_hub/internal/ip"
	"home_control_hub/internal/logrequest"
	"home_control_hub/internal/nas"
	"home_control_hub/internal/raspberry"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/", logrequest.LogRequest)

	// ip库
	ipGroup := r.Group("/ip")
	{
		ipGroup.GET("/", ip.QueryIp)
		ipGroup.GET("/update", ip.Update)
		ipGroup.GET("/:ip", ip.QueryIp)
	}
	// Nas
	nasGroup := r.Group("/nas")
	{
		nasGroup.GET("open", nas.Wake)
		nasGroup.GET("shutdown", nas.Shutdown)
	}
	// 树莓派
	raspGroup := r.Group("/raspberry")
	{
		raspGroup.GET("restart", raspberry.Restart)
		raspGroup.GET("shutdown", raspberry.Shutdown)
	}
	r.Run(":8082") // 监听并在 8082 上启动服务
}
