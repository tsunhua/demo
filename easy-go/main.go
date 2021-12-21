package main

import (
	"app/infrastructure/config"
	"app/infrastructure/debug"
	"app/infrastructure/log"
	"app/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type HttpMethod string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

type api struct {
	method  HttpMethod
	path    string
	handler gin.HandlerFunc
}

var apis = []api{
	{GET, "/ping", service.HandlePing},
}

func main() {
	debug.PprofRouter()
	router := gin.Default()
	for _, item := range apis {
		router.Handle(string(item.method), item.path, item.handler)
	}
	// 同源设置，node 调用后台接口需要
	c := cors.New(cors.Options{
		AllowedOrigins:   config.Get().CORS.AllowedOrigins,
		AllowCredentials: config.Get().CORS.AllowCredentials,
		// Enable Debugging for testing, consider disabling in production
		Debug: config.Get().CORS.Debug,
	})
	server := &http.Server{
		Addr:           ":" + config.Get().Port,
		Handler:        c.Handler(router),
		IdleTimeout:    6 * time.Minute,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("run server error", log.String("error", err.Error()))
	}
}
