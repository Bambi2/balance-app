package service

import (
	"database/sql"
	"net/http"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/repository"
)

type InvoiceService struct {
	repo repository.Invoice
}

func NewInvoiceRepository(repo repository.Invoice) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) Reserve(invoice domain.Invoice) error {
	hasMoney, err := s.repo.HasMoney(invoice.UserId, invoice.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return &CustomError{Message: ErrNoSuchUser, StatusCode: http.StatusBadRequest}
		}
		return err
	}

	if !hasMoney {
		return &CustomError{Message: ErrNotEnoughMoney, StatusCode: http.StatusBadRequest}
	}
	return s.repo.Reserve(invoice)
}

func (s *InvoiceService) Unreserve(userId int, serviceId int, orderId int64) error {
	hasUncheckedMoney, err := s.repo.HasUncheckedMoney(userId, serviceId, orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &CustomError{Message: ErrNoSuchInvoice, StatusCode: http.StatusBadRequest}
		}
		return err
	}

	if !hasUncheckedMoney {
		return &CustomError{Message: ErrCheckedInvoice, StatusCode: http.StatusBadRequest}
	}

	err = s.repo.Unreserve(userId, serviceId, orderId)
	if err == repository.ErrCheckedInvoice {
		return &CustomError{Message: ErrCheckedInvoice, StatusCode: http.StatusBadRequest}
	}

	return err
}
