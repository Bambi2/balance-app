package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	payForServiceDescription   = "%s куплена услуга: №%d, заказ: №%d | -%d RUB"
	unreservedMoneyDescription = "%s вернули деньги за услугу: №%d, заказ: №%d | +%d RUB "
)

var (
	ErrCheckedInvoice = errors.New("service was already applied")
)

type InvoicePostgres struct {
	db *sqlx.DB
}

func NewInvoicePostgres(db *sqlx.DB) *InvoicePostgres {
	return &InvoicePostgres{db: db}
}

func (r *InvoicePostgres) Reserve(invoice domain.Invoice) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createInvoiceQuery := fmt.Sprintf(`INSERT INTO %s (user_id, service_id, order_id, amount)
		VALUES($1, $2, $3, $4)`, invoicesTable)
	if _, err := tx.Exec(createInvoiceQuery, invoice.UserId, invoice.ServiceId, invoice.OrderId, invoice.Amount); err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now().Format(timeFormat)
	description := fmt.Sprintf(payForServiceDescription, now, invoice.ServiceId, invoice.OrderId, invoice.Amount/100)
	createNewTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)", transactionsTable)
	if _, err := tx.Exec(createNewTransactionQuery, invoice.UserId, invoice.Amount, description, now); err != nil {
		tx.Rollback()
		return err
	}

	takeMoneyFromUserQuery := fmt.Sprintf("UPDATE %s SET amount = amount - $1 WHERE id = $2", usersTable)
	if _, err = tx.Exec(takeMoneyFromUserQuery, invoice.Amount, invoice.UserId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *InvoicePostgres) HasMoney(userId int, price domain.CENT) (bool, error) {
	query := fmt.Sprintf("SELECT amount FROM %s WHERE id = $1", usersTable)

	var amount domain.CENT
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&amount); err != nil {
		return false, err
	}

	if amount < price {
		return false, nil
	}

	return true, nil
}

func (r *InvoicePostgres) Unreserve(userId int, serviceId int, orderId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var amountLeft domain.CENT
	deleteInvoice := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND service_id = $2 AND order_id = $3 RETURNING amount",
		invoicesTable)
	row := tx.QueryRow(deleteInvoice, userId, serviceId, orderId)
	if err := row.Scan(&amountLeft); err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now().Format(timeFormat)
	description := fmt.Sprintf(unreservedMoneyDescription, now, serviceId, orderId, amountLeft/100)
	createNewTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)",
		transactionsTable)
	if _, err := tx.Exec(createNewTransactionQuery, userId, amountLeft, description, now); err != nil {
		tx.Rollback()
		return err
	}

	if amountLeft == 0 {
		tx.Rollback()
		return ErrCheckedInvoice
	}

	returnMoneyQuery := fmt.Sprintf("UPDATE %s SET amount = amount + $1 WHERE id = $2", usersTable)
	if _, err := tx.Exec(returnMoneyQuery, amountLeft, userId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *InvoicePostgres) HasUncheckedMoney(userId int, serviceId int, orderId int64) (bool, error) {
	query := fmt.Sprintf("SELECT amount FROM %s WHERE user_id = $1 AND service_id = $2 AND order_id = $3", invoicesTable)

	var amount domain.CENT
	row := r.db.QueryRow(query, userId, serviceId, orderId)
	if err := row.Scan(&amount); err != nil {
		return false, err
	}

	if amount == 0 {
		return false, nil
	}

	return true, nil
}
