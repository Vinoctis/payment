package main 

import (
    "payment/internal"
	"os/signal"
	"log"
	"context"
	"os"
	"syscall"
	"time"
	"payment/pkg/config"
)

// var dbPayment *gorm.DB
var cfgManager *config.Manager

func init() {
	cfgManager = config.NewManager()
	//加载配置
	cfgManager.RegisterLoader(&config.EnvLoader{
		Envfile: ".env",
	})
	//数据库
	cfgManager.RegisterLoader(&config.YmalLoader{
		Path : "configs/database.yaml",
	})
	//注册中心
	cfgManager.RegisterLoader(&config.YmalLoader{
		Path : "configs/registry.yaml",
	})
	//服务
	cfgManager.RegisterLoader(&config.YmalLoader{
		Path : "configs/service.yaml",
	})

	if err := cfgManager.Load();err!= nil {
		panic(err)
	}
}

func main() {
	App, err := internal.NewApplication(cfgManager.Configs)
	if err != nil {
		panic(err)
	}
	
	// 开启服务
	go func() {
		if err := App.Server.HTTP.Run();err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	go func() {
		if err := App.Server.GRPC.Start();err!= nil {
			log.Printf("GRPC server error: %v", err)
		}
	}()
	log.Println("Server started")

	// 注册服务到etcd
	err = App.RegisterService()
	if err!= nil {
		log.Fatal(err)
	}
	
	//监听信号
	// - 注册要监听的系统信号到 quit 通道
	// - SIGINT : 当用户按下 Ctrl+C 时发送的中断信号
	// - SIGTERM : 程序终止信号，通常是优雅退出的信号
	quit := make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	//等待信号
	<-quit 
	log.Printf("Server shutting down...")

	ctx , cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//优雅关闭
	if err := App.Shutdown(ctx);err != nil {
		log.Printf("服务关闭出错:%v", err)
	}

	log.Println("服务已关闭")
}