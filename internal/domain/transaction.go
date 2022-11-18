package domain

type Transaction struct {
	Description string `json:"description" db:"description"`
}
