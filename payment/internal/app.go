package internal

import (
	"context"
	etcd "github.com/Vinoctis/etcd-register-center/sdk"
	"fmt"
	grpcHandler "payment/internal/handler/grpc"
	"payment/internal/handler/http/controller"
	"payment/internal/handler/http/router"
	repository "payment/internal/repository/mysql"
	"payment/internal/service"
	"payment/internal/transport/grpc"
	"payment/internal/transport/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"payment/pkg/config"
	mysql "payment/internal/repository/mysql"
)

type Application struct {
	config *config.Config
	databases *DBConnections
	providers *Providers 
	repositories *Repositories 
	services *Services
	controllers *Controllers
	Server *Servers
	router *gin.Engine 
	registry *Registry
}

type DBConnections struct {
	Payment *gorm.DB
}


type Providers struct {
	Payment map[string]service.PaymentProvider
}

type Repositories struct {
	Payment repository.PaymentRepository
}

type Services struct {
	Payment service.PaymentService
}

type Controllers struct {
	Payment *controller.PaymentController
}

type Servers struct {
	HTTP *http.Server
	GRPC *grpc.Server
}

type Registry struct {
	EtcdClient *etcd.Registry
	ServiceInfo map[string]*config.ServiceConfig
}



func NewApplication(cfg *config.Config) (*Application, error) {
	app := &Application{
		config: cfg,
		databases: &DBConnections{},
	}
	if err := app.InitComponents();err != nil {
		return nil , err
	}
	return app, nil
}

// 依赖注入管理
func (app *Application) InitComponents() error {

	if err := app.InitDatabases();err!= nil {
		return fmt.Errorf("数据库初始化失败")
	}

	if err := app.InitProviders();err!= nil {
		return fmt.Errorf("提供者初始化失败")
	}
	if err := app.InitRepositories();err != nil {
		return fmt.Errorf("仓库初始化失败")
	}
	if err := app.InitRepositories();err != nil {
		return fmt.Errorf("仓库初始化失败")
	}
	if err := app.InitServices();err != nil {
		return fmt.Errorf("服务初始化失败")
	}
	if err := app.InitControllers();err != nil {
		return fmt.Errorf("控制器初始化失败")
	}
	if err := app.InitRouter();err != nil {
		return fmt.Errorf("路由初始化失败")
	}
	if err := app.InitServer();err!= nil {
		return fmt.Errorf("服务初始化失败")
	}
	
	return nil
}

func (app *Application) InitDatabases() error {
		var err error
		app.databases.Payment ,err = mysql.Init(app.config.DB["payment"])
		if err!= nil {
			return fmt.Errorf("数据库初始化失败")
		}
		return nil
}

func (app *Application) InitProviders() error {
	app.providers = &Providers{
		Payment: map[string]service.PaymentProvider{
			//"wechat": service.NewWechatPaymentProvider(),
			"alipay": service.NewAlipayService(),
		},
	}
	return nil
}

func (app *Application) InitRepositories() error {
	app.repositories = &Repositories{
		Payment: repository.NewPaymentRepository(app.databases.Payment),
	}

	if app.repositories == nil {
		return fmt.Errorf("仓库初始化失败")
	}
	return nil
}

func (app *Application) InitServices() error {
	if app.repositories == nil {
		panic("仓库未初始化")
	}

    app.services = &Services{
        Payment: service.NewPaymentService(app.providers.Payment, app.repositories.Payment),
    }
    return nil
}

func (app *Application) InitControllers() error {
	if app.services == nil {
        return fmt.Errorf("服务未初始化")
    }
	app.controllers = &Controllers{
		Payment: controller.NewPaymentController(app.services.Payment),
	}
	return nil
}

func (app *Application) InitRouter() error {
	app.router = router.SetupRouter(app.controllers.Payment)
	return nil
}

func (app *Application) InitServer() error {
	httpAddr := fmt.Sprintf("%s:%s",
		app.config.Service["http"].Addr, app.config.Service["http"].Port)
	grpcAddr := fmt.Sprintf("%s:%s",
		app.config.Service["grpc"].Addr, app.config.Service["grpc"].Port)

	handler := grpcHandler.NewPaymentHanlder(app.services.Payment)
	app.Server = &Servers{
		HTTP: http.NewServer(httpAddr, app.router),
		GRPC: grpc.NewServer(grpcAddr, handler),
	}
	return nil
}

func (app *Application) RegisterService() error {
	etcdCfg := app.config.Registry
	//注册服务
	endPoints := []string{
		fmt.Sprintf("%s:%s", etcdCfg["etcd-node1"].Addr, etcdCfg["etcd-node1"].Port),
	} 
	client,err := etcd.NewEtcdClient(endPoints)
	if err!= nil {
		return fmt.Errorf("etcd client error: %v", err) 
	}

	app.registry = &Registry{
		EtcdClient: etcd.NewRegistry(client),
		ServiceInfo : app.config.Service,
	}
	
	for serviceType, info := range app.registry.ServiceInfo {
		serviceAddr := fmt.Sprintf("%s:%s", info.Addr, info.Port)
		if err := app.registry.EtcdClient.Register(info.Name,serviceAddr,5);err != nil {
			return fmt.Errorf("%v服务注册失败：%v", serviceType, err)
		}
		fmt.Printf("[%v]服务%v已注册，地址：%v\n", serviceType, info.Name, serviceAddr)
	}

	return nil
}

func (app *Application) DeregisterService() error {
	if app.registry != nil && app.registry.ServiceInfo != nil  {
		for serviceType, info := range app.registry.ServiceInfo {
			serviceAddr := fmt.Sprintf("%s:%s", info.Addr, info.Port)
			app.registry.EtcdClient.Deregister(info.Name, serviceAddr)
			fmt.Printf("[%v]服务已注销：%v\n", serviceType, info.Name)			
		}
		return nil
	} else {
		return fmt.Errorf("服务未注册")
	}
}

func (app *Application) Shutdown(ctx context.Context) error {
	//注销服务
	if err := app.DeregisterService();err!= nil {
		return err
	}
	//关闭网络服务
	if err := app.Server.HTTP.Shutdown(ctx);err!= nil {
		return err
	}
	err := app.Server.GRPC.Stop(ctx)
	if err!= nil {
		return err
	}

	//关闭数据库链接
	sqlDB, err := app.databases.Payment.DB()
	if err != nil {
		return err 
	}
	return sqlDB.Close()
}
