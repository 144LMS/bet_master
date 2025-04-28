package models

import (
	"time"
)

type Wallet struct {
	ID           uint          `gorm:"primaryKey"`
	UserID       uint          `gorm:"not null;uniqueIndex"`
	Balance      float64       `gorm:"not null;default:0"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:WalletID"`
}

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	WalletID    uint      `gorm:"not null;index"` // Индекс для быстрого поиска
	Amount      float64   `gorm:"not null"`
	Type        string    `gorm:"not null;type:varchar(20)"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index"` // Индекс для сортировки
}
