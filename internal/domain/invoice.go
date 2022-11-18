package domain

type Invoice struct {
	UserId    int   `json:"userId" binding:"required"`
	ServiceId int   `json:"serviceId" binding:"required"`
	OrderId   int64 `json:"orderId" binding:"required"`
	Amount    CENT  `json:"amount" binding:"required"`
}
