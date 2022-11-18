package domain

type Check struct {
	UserId    int   `json:"userId" binding:"required"`
	ServiceId int   `json:"serviceId" binding:"required"`
	OrderId   int64 `json:"orderId" binding:"required"`
	Amount    CENT  `json:"amount" binding:"required"`
}

type CheckForRecords struct {
	ServiceId int  `db:"service_id"`
	Amount    CENT `db:"amount"`
}
