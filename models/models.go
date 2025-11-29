package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:500;not null" json:"name"`
	Email        string    `gorm:"size:500;uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"size:500;not null" json:"-"`
	RefreshToken string    `gorm:"type:text" json:"-"`
	AccessToken  string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Account struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BankName  string    `gorm:"size:255;not null" json:"bank_name"`
	Amount    float64   `gorm:"type:decimal(12,2);not null;default:0" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	CategoryID *uint     `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	AccountID  *uint     `gorm:"index" json:"account_id"`
	Account    *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type TransactionResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Amount       float64   `json:"amount"`
	CategoryID   *uint     `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	AccountID    *uint     `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Budget struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Criteria   string    `gorm:"size:50;not null" json:"criteria"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type ScheduledTransaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Repetition string    `gorm:"size:50;not null" json:"repetition"`
	RepeatAt   time.Time `gorm:"not null" json:"repeat_at"`
	CategoryID *uint     `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	AccountID  *uint     `gorm:"index" json:"account_id"`
	Account    *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type CategoryCreate struct {
	Name string `json:"name" binding:"required"`
}
type TransactionCreate struct {
	Name       string  `json:"name" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	CategoryID *uint   `json:"category_id"`
	AccountID  *uint   `json:"account_id"`
}

type BudgetCreate struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Criteria   string  `json:"criteria" binding:"required"`
}

type BudgetResponse struct {
	ID           uint      `json:"id"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	Amount       float64   `json:"amount"`
	Criteria     string    `json:"criteria"`
	Spent        float64   `json:"spent"`
	Remaining    float64   `json:"remaining"`
	Percentage   float64   `json:"percentage"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ScheduledTransactionCreate struct {
	Name       string    `json:"name" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Repetition string    `json:"repetition" binding:"required"`
	RepeatAt   time.Time `json:"repeat_at" binding:"required"`
	CategoryID *uint     `json:"category_id"`
	AccountID  *uint     `json:"account_id"`
}

type ScheduledTransactionResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Amount       float64   `json:"amount"`
	Repetition   string    `json:"repetition"`
	RepeatAt     time.Time `json:"repeat_at"`
	CategoryID   *uint     `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	AccountID    *uint     `json:"account_id"`
	AccountName  string    `json:"account_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
