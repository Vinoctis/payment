 package grpc 

 import (
	pb "payment/api/proto"
	"payment/internal/service"
	"payment/internal/validators"
	//"gin-gonic/gin"
	"context"
 )

 type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
 }

 func NewPaymentHanlder(ps service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: ps,
	}
 }

 func (h *PaymentHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	serviceReq := &validators.PaymentCreateOrderReq{
		UserId: int(req.UserId),
		Amount: int(req.Amount),
		Payway: req.Payway,
	}
	//service 方法
	order, err := h.paymentService.CreateOrder(serviceReq)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		OrderNo: order.OrderNo,
		Status: int32(order.Status),
	}, nil
 }

 func (h *PaymentHandler) HandleNotify(ctx context.Context, req *pb.NotifyRequest) (*pb.NotifyResponse, error) {
	orderNo := req.OrderNo
	data := req.Data
	result, err := h.paymentService.HandleNotify(orderNo, data)
	if err!= nil {
		return nil, err
	}
	return &pb.NotifyResponse{
		OrderNo: result.OrderNo,
		Status: int32(result.Status),
	}, nil
 }
