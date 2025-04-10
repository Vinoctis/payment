package controller

import (
	"payment/internal/service"
	"github.com/gin-gonic/gin"
	"payment/internal/validators"
	"payment/pkg/utils"
	"payment/internal/model"
)

type PaymentController struct {
	service service.PaymentService
}

func (pc *PaymentController) isValidPayway(name string) (*model.Payway, bool) {
	payway, err := pc.service.ValidatePayway(name)
	if err!= nil {
		return nil, false
	}
	return payway , true
}

func NewPaymentController(service service.PaymentService) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

func (pc *PaymentController) CreateOrder(c *gin.Context) {
	payway := c.Param("payway")
	paywayRepo, ok := pc.isValidPayway(payway)
	if !ok {
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, "无效的支付方式")
		return
	}

	var req validators.PaymentCreateOrderReq
	req.UserId = c.GetInt("user_id")
	req.PaywayId = int(paywayRepo.ID)
	if err := c.ShouldBind(&req);err != nil {
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, err.Error())
		return
	}

 	order, err := pc.service.CreateOrder(&req)
	if err!= nil {
		utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, err.Error())
		return
	}
	utils.ResponseSuccess(c, order)
}

func (pc *PaymentController) HandleNotify(c *gin.Context)  {
	orderNo := c.Query("order_no")
	data, err := c.GetRawData()
	if err != nil {
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParam, err.Error())
		return
	}
	result, err := pc.service.HandleNotify(orderNo, data)
	if err!= nil {
		utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, err.Error())
		return 
	}
	utils.ResponseSuccess(c, result)
}

