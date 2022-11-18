package service

import (
	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/repository"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactions(userId int, limit int, offset int) ([]domain.Transaction, error) {
	return s.repo.GetTransactions(userId, limit, offset)
}
