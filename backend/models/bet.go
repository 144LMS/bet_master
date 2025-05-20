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
	UserID       uint      `gorm:"not null;index"`
	WalletID     uint      `gorm:"not null;index"`
	MatchID      uint      `gorm:"not null;index"`
	Amount       float64   `gorm:"not null"`
	Odds         float64   `gorm:"not null"`
	PotentialWin float64   `gorm:"not null"`
	Team         string    `gorm:"not null;size:10"`
	Status       BetStatus `gorm:"default:'pending'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	SettledAt    *time.Time
}
