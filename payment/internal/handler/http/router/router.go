package router

import (
	"payment/internal/handler/http/controller"
	"payment/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(paymentController *controller.PaymentController) *gin.Engine {
	router := gin.Default()
	
	//下单
	router.POST("/order/:payway", paymentController.CreateOrder).Use(middleware.AuthMiddleware())
	router.POST("/notify", paymentController.HandleNotify)
	
	return router
}