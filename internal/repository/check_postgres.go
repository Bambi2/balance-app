package repository

import (
	"fmt"
	"time"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CheckPostgres struct {
	db *sqlx.DB
}

func NewCheckPostgres(db *sqlx.DB) *CheckPostgres {
	return &CheckPostgres{db: db}
}

func (r *CheckPostgres) Check(check domain.Check) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createCheckQuery := fmt.Sprintf("INSERT INTO %s (service_id, amount, created_at) VALUES($1, $2, $3)", checksTable)
	if _, err := tx.Exec(createCheckQuery, check.ServiceId, check.Amount, time.Now().Format(timeFormat)); err != nil {
		tx.Rollback()
		return err
	}

	takeMoneyFromInvoiceQuery := fmt.Sprintf("UPDATE %s SET amount = amount - $1 WHERE user_id = $2 AND service_id = $3 AND order_id = $4", invoicesTable)
	if _, err := tx.Exec(takeMoneyFromInvoiceQuery, check.Amount, check.UserId, check.ServiceId, check.OrderId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *CheckPostgres) HasMoney(check domain.Check) (bool, error) {
	query := fmt.Sprintf("SELECT amount FROM %s WHERE service_id = $1 AND order_id = $2 AND user_id = $3", invoicesTable)

	var amount domain.CENT
	row := r.db.QueryRow(query, check.ServiceId, check.OrderId, check.UserId)
	if err := row.Scan(&amount); err != nil {
		return false, err
	}

	if amount < check.Amount {
		return false, nil
	}

	return true, nil
}

func (r *CheckPostgres) GetChecks(startDate string, endDate string) ([]domain.CheckForRecords, error) {
	query := fmt.Sprintf("SELECT service_id, amount FROM %s WHERE created_at >= $1 and created_at < $2", checksTable)

	var checks []domain.CheckForRecords
	err := r.db.Select(&checks, query, startDate, endDate)

	return checks, err
}
