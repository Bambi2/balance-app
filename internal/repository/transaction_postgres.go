package repository

import (
	"fmt"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) GetTransactions(userId int, limit int, offset int) ([]domain.Transaction, error) {
	query := fmt.Sprintf("SELECT description FROM %s WHERE user_id = $1 ORDER BY created_at DESC, amount DESC  LIMIT $2 OFFSET $3", transactionsTable)

	var transactions []domain.Transaction
	err := r.db.Select(&transactions, query, userId, limit, offset)

	return transactions, err
}
