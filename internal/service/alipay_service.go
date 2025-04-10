package service 

import (
	"payment/internal/model"
)

type AlipayService struct {}

func NewAlipayService() PaymentProvider {
	return &AlipayService{}
}

func (s *AlipayService) GeneratePayURL(order *model.Order) (string, error) {
	return "http://alipay_url", nil
} 

func (s *AlipayService) VerifyNotify(data []byte) bool {
	return true
} 