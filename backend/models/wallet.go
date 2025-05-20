package models

import (
	"time"
)

type Wallet struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	UserID       uint          `json:"user_id" gorm:"not null;uniqueIndex"`
	Balance      float64       `json:"balance" gorm:"not null;default:0"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:WalletID"`
}

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WalletID    uint      `json:"wallet_id" gorm:"not null;index"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null;type:varchar(20)"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}
