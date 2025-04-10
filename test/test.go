package main 

import (
	"fmt"
	// "payment/pkg/utils"
	repository "payment/internal/repository"
	mysql "payment/internal/repository/mysql"
	"payment/internal/service"
	"etcd-register-center/sdk"
	"log"

)
func test3() {
	defer fmt.Println("before panic")
	panic("this is panic ")
	defer fmt.Println("after panic")
}

func test2() int {
	x := 1
	defer func(){
		x++
		fmt.Println("defer:",x)
	}()
	return x+2
}
func main() {
	
	DBCfg,err := repository.LoadConfig()
	if err != nil {
		panic(err)
	}
	//初始化数据库
	db , err := mysql.Init(DBCfg.DBPayment.MySqlDB)
	if err != nil {
		panic(err)
	}
	repo := mysql.NewPaymentRepository(db)
	provider := map[string]service.PaymentProvider{
		"alipay" : service.NewAlipayService(),
	}
	paymentService := service.NewPaymentService(provider, repo)
	fmt.Println(paymentService)

	etcdEndpoints := []string{"http://localhost:2379"}
	serviceName := "payment"
	address := "127.0.0.1:9090"

	// 创建 etcd 客户端
	client, err := sdk.NewEtcdClient(etcdEndpoints)
	fmt.Println(client)
	if err != nil {
		panic(err)
	}
	registry := sdk.NewRegistery(client)
	// 注册服务
	err = registry.Register(serviceName, address, 5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Service registered:", serviceName, address)

}