package validators

type PaymentCreateOrderReq struct {
	UserId int `json:"user_id" form:"user_id" binding:"required"`
	Payway string `json:"payway" form:"payway" binding:"required"`
	PaywayId int `json:"payway_id" form:"payway_id" binding:"omitempt":` 
	Amount int `json:"amount" form:"amount" binding:"required"`
	GameId int `json:"game_id" form:"game_id" binding:"required"`
}