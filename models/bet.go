package models

import (
	"time"

	"gorm.io/gorm"
)

type BetStatus string

const (
	BetPending BetStatus = "pending"
	BetWon     BetStatus = "won"
	BetLost    BetStatus = "lost"
)

type Bet struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	Amount       float64   `gorm:"not null"`
	Odds         float64   `gorm:"not null"` // Коэффициент (например, 1.95)
	PotentialWin float64   `gorm:"not null"` // amount * odds
	Status       BetStatus `gorm:"default:'pending'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	SettledAt    *time.Time
}
