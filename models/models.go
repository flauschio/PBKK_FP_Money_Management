package models

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	CategoryID   *int      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	Amount       float64   `json:"amount"`
	AccountID    *int      `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TransactionCreate struct {
	Name       string  `json:"name"`
	CategoryID *int    `json:"category_id"`
	Amount     float64 `json:"amount"`
	AccountID  *int    `json:"account_id"`
}

type CategoryCreate struct {
	Name string `json:"name"`
}

// hardcoded data for future features
type Account struct {
	ID        int       `json:"id"`
	BankName  string    `json:"bank_name"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Budget struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Amount     float64   `json:"amount"`
	Criteria   string    `json:"criteria"`
	UpdatedAt  time.Time `json:"updated_at"`
}
