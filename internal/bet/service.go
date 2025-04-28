package bet

import (
	"time"

	"github.com/144LMS/bet_master/models"
)

type BetService struct {
	betRepo BetRepository
}

func NewBetService(betRepo BetRepository) *BetService {
	return &BetService{
		betRepo: betRepo,
	}
}

func (s *BetService) PlaceBet(bet models.Bet) (*models.Bet, error) {
	return s.betRepo.CreateBet(&bet)
}

func (s *BetService) SettleBet(betID uint, won bool) (*models.Bet, error) {
	bet, err := s.betRepo.GetBet(betID)
	if err != nil {
		return nil, err
	}

	status := models.BetLost
	if won {
		status = models.BetWon
	}

	bet.Status = status
	now := time.Now()
	bet.SettledAt = &now

	return s.betRepo.UpdateBet(bet)
}
