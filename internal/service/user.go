package service

import (
	"net/http"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddToBalance(user domain.User) (domain.CENT, error) {
	if user.Amount < 0 {
		return 0, &CustomError{Message: ErrNegativeAmount, StatusCode: http.StatusBadRequest}
	}
	return s.repo.AddToBalance(user)
}

func (s *UserService) GetBalance(userId int) (domain.CENT, error) {
	userExists, err := s.repo.IfExists(domain.User{Id: userId})
	if err != nil {
		return 0, err
	}

	if !userExists {
		return 0, &CustomError{Message: ErrNoSuchUser, StatusCode: http.StatusBadRequest}
	}

	return s.repo.GetBalance(userId)
}

func (s *UserService) SendMoney(senderId int, recieverId int, amount domain.CENT) error {
	if amount <= 0 {
		return &CustomError{Message: "amount cannot be less or equal to 0", StatusCode: http.StatusBadRequest}
	}
	senderExists, err := s.repo.IfExists(domain.User{Id: senderId})
	if err != nil {
		return err
	}

	if !senderExists {
		return &CustomError{Message: ErrSenderNotExists, StatusCode: http.StatusBadRequest}
	}

	recieverExists, err := s.repo.IfExists(domain.User{Id: recieverId})
	if err != nil {
		return err
	}

	if !recieverExists {
		return &CustomError{Message: ErrReciverNotExists, StatusCode: http.StatusBadRequest}
	}

	senderBalance, err := s.repo.GetBalance(senderId)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return &CustomError{Message: ErrSenderNotEnoughMoney, StatusCode: http.StatusBadRequest}
	}

	return s.repo.SendMoney(senderId, recieverId, amount)
}
