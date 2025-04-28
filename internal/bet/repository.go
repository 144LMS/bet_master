package bet

import (
	"github.com/144LMS/bet_master/models"
	"gorm.io/gorm"
)

type BetRepository struct {
	db *gorm.DB
}

func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{db: db}
}

func (r *BetRepository) CreateBet(bet *models.Bet) (*models.Bet, error) {
	err := r.db.Create(bet).Error
	return bet, err
}

func (r *BetRepository) GetBet(id uint) (*models.Bet, error) {
	var bet models.Bet
	err := r.db.First(&bet, id).Error
	return &bet, err
}

func (r *BetRepository) UpdateBet(bet *models.Bet) (*models.Bet, error) {
	err := r.db.Save(bet).Error
	return bet, err
}
