package models

import "gorm.io/gorm"

type Bet struct {
	gorm.Model
	FirstTeam  string  `json:"first_team"`
	SecondTeam string  `json:"second_team"`
	Winner     string  `json:"winner"`
	Total      float64 `json:"total"`
}

type UpdateBetRequest struct {
	IDBet      int     `json:"id_bet"`
	FirstTeam  string  `json:"first_team"`
	SecondTeam string  `json:"second_team"`
	Total      float64 `json:"total"`
}

type WinnerBetRequest struct {
	IDBet  int    `json:"id_bet"`
	Winner string `json:"winner"`
}

type DeleteBetRequest struct {
	IDBet int `json:"id_bet"`
}
