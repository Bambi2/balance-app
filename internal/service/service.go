package service

import (
	"os"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/repository"
)

type CustomError struct {
	Message    string
	StatusCode int
}

func (e *CustomError) Error() string {
	return e.Message
}

var (
	ErrNegativeAmount       = "cannot add negative values"
	ErrNoSuchUser           = "user with such id doesn't exist"
	ErrSenderNotExists      = "user with such senderId doesn't exist"
	ErrReciverNotExists     = "user with such recieverId doesn't exist"
	ErrSenderNotEnoughMoney = "sender doesn't have such amount"
	ErrNotEnoughMoney       = "user doesn't have enough money for the service"
	ErrNoSuchInvoice        = "such invoice doesn't exist"
	ErrNoMoneyLeftInvoice   = "no money left in invoice"
	ErrCheckedInvoice       = "the invoice is already checked"
)

type User interface {
	AddToBalance(user domain.User) (domain.CENT, error)
	GetBalance(userId int) (domain.CENT, error)
	// may return ErrSenderNotExists, ErrReciverNotExists, ErrSenderNotEnoughMoney
	SendMoney(senderId int, recieverId int, amount domain.CENT) error
}

type Invoice interface {
	Reserve(invoice domain.Invoice) error
	Unreserve(userId int, serviceId int, orderId int64) error
}

type Transaction interface {
	GetTransactions(userId int, limit int, offset int) ([]domain.Transaction, error)
}

type Check interface {
	Check(check domain.Check) error
	//it is caller's responsobility to remove the returned file
	GetCSVChecks(year string, month string) (*os.File, error)
}

type Service struct {
	User
	Invoice
	Transaction
	Check
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:        NewUserService(repo.User),
		Invoice:     NewInvoiceRepository(repo.Invoice),
		Check:       NewCheckService(repo.Check),
		Transaction: NewTransactionService(repo.Transaction),
	}
}
