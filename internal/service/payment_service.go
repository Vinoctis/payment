package service 

import (
	"payment/internal/model"
	repository "payment/internal/repository/mysql"
	"payment/pkg/utils"
	"payment/internal/validators"
	"errors"
	"fmt"
)

var (
	ErrInvalidPayway = errors.New("支付方式不合法")
	ErrVerifyNotify = errors.New("验证失败")
)

type NotifyResult struct {
	OrderNo string
	Status int32
	PlatformOrderNo string
	Subject string
}

type PaymentService interface { 
	CreateOrder(req *validators.PaymentCreateOrderReq) (*model.Order, error) //下单
	HandleNotify(ordreNo string, data []byte) (*NotifyResult, error) //notify
	ValidatePayway(name string) (*model.Payway, error)
}

type PaymentProvider interface {
	GeneratePayURL(order *model.Order) (string, error)
	VerifyNotify(data []byte) bool
}

type paymentService struct {
	provider map[string]PaymentProvider
	repo repository.PaymentRepository
}

func NewPaymentService(provider map[string]PaymentProvider, repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		// provider: map[string]PaymentProvider{
		// 	"alipay": NewAlipayService(),
		// 	//"wechat": NewWechatService(),
		// },
		provider: provider,
		repo: repo,
	}
}

func (s *paymentService) CreateOrder(req *validators.PaymentCreateOrderReq) (*model.Order, error) {
	//获取支付方式
	provider, ok := s.provider[req.Payway]
	if !ok {
		return nil, ErrInvalidPayway
	}

	orderNo := utils.GenerateOrderID()
	order := &model.Order{
		Amount: req.Amount,
		PaywayId: req.PaywayId,
		OrderNo: orderNo,
		Status: model.OrderStatusPending,
	}
	err := s.repo.Order().Create(order)
	if err != nil {
		return nil, err
	}
	payURL, err := provider.GeneratePayURL(order)
	if err != nil {
		return nil, err
	}
	//
	fmt.Println("payURL:", payURL)
	//请求支付网关
	//...
	err = s.repo.Order().Update(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *paymentService) HandleNotify(orderNo string, data []byte) (*NotifyResult, error) {
	order , err := s.repo.Order(). FindByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	payway,err := s.repo.Order().GetPayway(order)
	if err != nil{
		return nil, ErrInvalidPayway
	}
	provider, ok := s.provider[payway.Name]
	if!ok {
		return nil, ErrInvalidPayway
	}

	if provider.VerifyNotify(data) {
		//更新订单状态
		order.Status = model.OrderStatusPaid 
		err := s.repo.Order().Update(order)
		if err != nil {
			return nil, err
		}
		return &NotifyResult{
			OrderNo: order.OrderNo,
			Status: order.Status,
			//PlatformOrderNo: order.PlatformOrderNo,
		}, nil
	} else {
		return nil, ErrVerifyNotify
	}
}

func (s *paymentService) ValidatePayway(name string) (*model.Payway, error) {
	payway, err := s.repo.Payway().FindByName(name)
	if err != nil {
		return nil, err
	}
	return payway, nil
}

