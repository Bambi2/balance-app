package domain

type User struct {
	Id     int  `json:"userId" binding:"required"`
	Amount CENT `json:"amount" binding:"required"`
}
