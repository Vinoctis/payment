package internal

import (
	"context"
	etcd "etcd-register-center/sdk"
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
)

type Application struct {
	DB *gorm.DB
	providers *Providers 
	repositories *Repositories 
	services *Services
	controllers *Controllers
	Server *Servers
	router *gin.Engine 
	registry *Registry
	etcdServiceInfo *EtcdServiceInfo
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

type EtcdServiceInfo struct {
	Name string
	Addr string
}

type Registry struct {
	EtcdClient *etcd.Registry
	ServiceInfo map[string]*EtcdServiceInfo 
}



func NewApplication(db *gorm.DB) (*Application, error) {
	app := &Application{
		DB: db,
	}
	if err := app.InitComponents();err != nil {
		return nil , err
	}
	return app, nil
}

// 依赖注入管理
func (app *Application) InitComponents() error {

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
		Payment: repository.NewPaymentRepository(app.DB),
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
	handler := grpcHandler.NewPaymentHanlder(app.services.Payment)
	app.Server = &Servers{
		HTTP: http.NewServer(":8080", app.router),
		GRPC: grpc.NewServer(":9090", handler),
	}
	return nil
}

func (app *Application) RegisterService() error {
	//注册服务
	endPoints := []string{
		"localhost:2379",
	} 
	client,err := etcd.NewEtcdClient(endPoints)
	if err!= nil {
		return fmt.Errorf("etcd client error: %v", err) 
	}

	info := make(map[string]*EtcdServiceInfo)
	info["HTTP"] = &EtcdServiceInfo{
		Name: "payment-http",
		Addr: "localhost:8080",
	}
	info["GRPC"] = &EtcdServiceInfo{
		Name: "payment-grpc",
		Addr: "localhost:9090",
	}

	app.registry = &Registry{
		EtcdClient: etcd.NewRegistry(client),
		ServiceInfo : info,
	}
	
	for serviceType, info := range app.registry.ServiceInfo {
		if err := app.registry.EtcdClient.Register(info.Name,info.Addr,5);err != nil {
			return fmt.Errorf("%v服务注册失败：%v", serviceType, err)
		}
		fmt.Printf("[%v]服务%v已注册，地址：%v\n", serviceType, info.Name, info.Addr)
	}

	return nil
}

func (app *Application) DeregisterService() error {
	if app.registry != nil && app.registry.ServiceInfo != nil  {
		for serviceType, info := range app.registry.ServiceInfo {
			app.registry.EtcdClient.Deregister(info.Name, info.Addr)
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
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err 
	}
	return sqlDB.Close()
}
