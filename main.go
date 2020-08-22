package main

import (
	"context"
	"dj-api/app/api"
	"dj-api/app/config"
	"dj-api/app/rpc/service"
	"dj-api/dal/db"
	"dj-api/dal/grpc/server"
	"dj-api/dal/nsq"
	"dj-api/dal/redis"
	"dj-api/tools"
	"dj-api/tools/logger"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := config.SetUp()
	if err != nil {
		fmt.Printf("config setup err:%v\n", err)
		panic(err)
	}
	err = db.Setup()
	if err != nil {
		fmt.Printf("db setup err:%v\n", err)
		panic(err)
	}
	redis.Setup()
	err = logger.Setup(config.C.Registry.ServiceName, config.C.Log.Dir, logger.LogLevel(config.C.Log.Level))
	if err != nil {
		fmt.Printf("log setup err:%v\n", err)
		panic(err)
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.C.Sentry.Dsn,
	}); err != nil {
		logger.Panic("sentry init err:%v", err)
	}
}

func main() {
	gin.SetMode(config.C.Api.RunMode)

	//日志
	if tools.IsDev() {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		// 创建记录日志的文件
		f, err := os.Create(tools.GetRootDir() + "/logs/gin.log")
		if err != nil {
			logger.Panic("gin log file create fatal:%v", err)
		}
		gin.DefaultWriter = io.MultiWriter(f)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.C.Api.Port),
		Handler:        api.InitRouter(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic("http server start fatal:%v", err)
		}
	}()

	//开启nsq消费者
	err := nsq.Setup()
	if err != nil {
		logger.Panic("nsq consumer start fatal:%v", err)
	}
	//开启gRpc server
	err = server.Start()
	if err != nil {
		logger.Panic("server start fatal:%v", err)
	}
	//开启gRpc client
	err = service.Init()
	if err != nil {
		logger.Panic("client start fatal:%v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Server Https Shutdown:", err)
	}
	shutdown()
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		fmt.Println("timeout of 2 seconds.")
	}
	fmt.Println("Server exiting success!!!")
}

func shutdown() {
	err := db.CloseDB()
	if err != nil {
		logger.Panic("db close err:%v", err)
	}
	err = redis.Close()
	if err != nil {
		logger.Panic("redis close err:%v", err)
	}
	nsq.StopC()
	server.Stop()
	err = service.Close()
	if err != nil {
		logger.Panic("gRpc client close:%v", err)
	}
	sentry.Flush(time.Second)
}
