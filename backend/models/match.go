package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchUpcoming  MatchStatus = "upcoming"
	MatchLive      MatchStatus = "live"
	MatchCompleted MatchStatus = "completed"
	MatchCancelled MatchStatus = "cancelled"
)

type Match struct {
	gorm.Model
	Team1     string      `gorm:"not null;size:100" json:"team1"`
	Team2     string      `gorm:"not null;size:100" json:"team2"`
	StartTime time.Time   `gorm:"not null;index" json:"start_time"`
	EndTime   *time.Time  `json:"end_time,omitempty"`
	Status    MatchStatus `gorm:"type:varchar(20);default:'upcoming'" json:"status"`
	SportType string      `gorm:"size:50" json:"sport_type"`
}
