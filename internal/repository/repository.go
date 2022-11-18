package repository

import (
	"github.com/bambi2/balance-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type User interface {
	AddToBalance(user domain.User) (domain.CENT, error)
	GetBalance(userId int) (domain.CENT, error)
	IfExists(user domain.User) (bool, error)
	SendMoney(senderId int, recieverId int, amount domain.CENT) error
}

type Invoice interface {
	Reserve(invoice domain.Invoice) error
	HasMoney(userId int, amount domain.CENT) (bool, error)
	Unreserve(userId int, serviceId int, orderId int64) error
	HasUncheckedMoney(userId int, serviceId int, orderId int64) (bool, error)
}

type Transaction interface {
	GetTransactions(userId int, limit int, offset int) ([]domain.Transaction, error)
}

type Check interface {
	Check(check domain.Check) error
	HasMoney(check domain.Check) (bool, error)
	GetChecks(startDate string, endDate string) ([]domain.CheckForRecords, error)
}

type Repository struct {
	User
	Invoice
	Transaction
	Check
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserPostgres(db),
		Invoice:     NewInvoicePostgres(db),
		Check:       NewCheckPostgres(db),
		Transaction: NewTransactionPostgres(db),
	}
}
