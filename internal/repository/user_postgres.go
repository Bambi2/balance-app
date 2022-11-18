package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	addToBalanceDescription                = "%s начислено на баланс | +%d RUB"
	sendMoneyToAnotherUserDescription      = "%s перечислены деньги на баланс пользователя: %d | -%d RUB"
	recieveMoneyFromAnotherUserDescription = "%s получены деньги от пользователя: %d | +%d RUB"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetBalance(userId int) (domain.CENT, error) {
	query := fmt.Sprintf("SELECT amount FROM %s WHERE id = $1", usersTable)

	var amount domain.CENT
	row := r.db.QueryRow(query, userId)
	err := row.Scan(&amount)

	return amount, err
}

func (r *UserPostgres) AddToBalance(user domain.User) (domain.CENT, error) {
	userExists, err := r.IfExists(user)
	if err != nil {
		return 0, err
	}

	if !userExists {
		return r.createNewBalance(user)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var newAmount domain.CENT
	addToBalanceQuery := fmt.Sprintf("UPDATE %s SET amount = amount + $1 WHERE id = $2 RETURNING amount", usersTable)
	row := tx.QueryRow(addToBalanceQuery, user.Amount, user.Id)
	if err = row.Scan(&newAmount); err != nil {
		tx.Rollback()
		return 0, err
	}

	now := time.Now().Format(timeFormat)
	description := fmt.Sprintf(addToBalanceDescription, now, user.Amount/100)
	createNewTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)", transactionsTable)
	if _, err = tx.Exec(createNewTransactionQuery, user.Id, user.Amount, description, now); err != nil {
		tx.Rollback()
		return 0, err
	}

	return newAmount, tx.Commit()
}

func (r *UserPostgres) SendMoney(senderId int, recieverId int, amount domain.CENT) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Format(timeFormat)

	senderDescription := fmt.Sprintf(sendMoneyToAnotherUserDescription, now, recieverId, amount/100)
	createNewSenderTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)",
		transactionsTable)
	if _, err := tx.Exec(createNewSenderTransactionQuery, senderId, amount, senderDescription, now); err != nil {
		tx.Rollback()
		return err
	}

	recieverDescription := fmt.Sprintf(recieveMoneyFromAnotherUserDescription, now, senderId, amount/100)
	createNewRecieverTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)",
		transactionsTable)
	if _, err := tx.Exec(createNewRecieverTransactionQuery, senderId, amount, recieverDescription, now); err != nil {
		tx.Rollback()
		return err
	}

	addToRecieverBalaneQuery := fmt.Sprintf("UPDATE %s SET amount = amount + $1 WHERE id = $2", usersTable)
	if _, err := tx.Exec(addToRecieverBalaneQuery, amount, recieverId); err != nil {
		tx.Rollback()
		return err
	}

	decreaseSendersBalaneQuery := fmt.Sprintf("UPDATE %s SET amount = amount + $1 WHERE id = $2", usersTable)
	if _, err := tx.Exec(decreaseSendersBalaneQuery, amount, senderId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserPostgres) createNewBalance(user domain.User) (domain.CENT, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createNewBalanceQuery := fmt.Sprintf("INSERT INTO %s (id, amount) VALUES ($1, $2)", usersTable)
	if _, err = tx.Exec(createNewBalanceQuery, user.Id, user.Amount); err != nil {
		tx.Rollback()
		return 0, err
	}

	now := time.Now().Format(timeFormat)
	description := fmt.Sprintf(addToBalanceDescription, now, user.Amount/100)
	createNewTransactionQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4)", transactionsTable)
	if _, err = tx.Exec(createNewTransactionQuery, user.Id, user.Amount, description, now); err != nil {
		tx.Rollback()
		return 0, err
	}

	return user.Amount, tx.Commit()
}

func (r *UserPostgres) IfExists(user domain.User) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT exists(SELECT * FROM %s WHERE id = $1)", usersTable)
	err := r.db.QueryRow(query, user.Id).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}
