package main

import "github.com/gin-gonic/gin"

func createRouter(r *gin.RouterGroup) {
	r.Any("/", help)
	r.Any("/help", help)
	r.Any("/ping", ping)
	r.Any("/ready", isReady)
	r.Any("/info", info)
	r.Any("/time", timeData)
	r.Any("/echo", echo)
	r.Any("/crash", crash)
	r.Any("/wait/:ms", wait)
	r.Any("/exit", exit)
	r.Any("/status/:code", statusCode)
	r.GET("/file", readFile)
	r.POST("/file", uploadFile)
}
