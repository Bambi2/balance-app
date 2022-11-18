package service

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/repository"
)

type CheckService struct {
	repo repository.Check
}

func NewCheckService(repo repository.Check) *CheckService {
	return &CheckService{repo: repo}
}

func (s *CheckService) Check(check domain.Check) error {
	enoughMoney, err := s.repo.HasMoney(check)
	if err != nil {
		if err == sql.ErrNoRows {
			return &CustomError{Message: ErrNoSuchInvoice, StatusCode: http.StatusBadRequest}
		}
		return err
	}

	if !enoughMoney {
		return &CustomError{Message: ErrNoMoneyLeftInvoice, StatusCode: http.StatusBadRequest}
	}

	return s.repo.Check(check)
}

func (s *CheckService) GetCSVChecks(year string, month string) (*os.File, error) {
	// formating date for repository
	if year == "" || month == "" {
		return nil, &CustomError{Message: "empty data", StatusCode: http.StatusBadRequest}
	}

	startDate, err := time.Parse("2006-01", year+"-"+month)
	if err != nil {
		return nil, &CustomError{Message: "invalid input date", StatusCode: http.StatusBadRequest}
	}
	endDate := startDate.AddDate(0, 1, 0)

	// getting checks from repository
	checks, err := s.repo.GetChecks(startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	if checks == nil {
		return nil, &CustomError{Message: "no data for this period", StatusCode: http.StatusBadRequest}
	}

	// sorting checks by ids
	checksSorted := make(map[int]domain.CENT)
	for _, val := range checks {
		checksSorted[val.ServiceId] += val.Amount
	}

	// formating data for csv file
	var finalChecks [][]string = make([][]string, 0, len(checksSorted))
	for key, val := range checksSorted {
		finalChecks = append(finalChecks, []string{strconv.Itoa(key), strconv.Itoa(int(val))})
	}

	// creating temp csv file
	tempFile, err := os.CreateTemp("", "data.*.csv")
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(tempFile)

	csvWriter.WriteAll(finalChecks)

	return tempFile, nil
}
