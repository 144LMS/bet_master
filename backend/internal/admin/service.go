package admin

import (
	"errors"
	"time"

	"github.com/144LMS/bet_master/internal/wallet"
	"github.com/144LMS/bet_master/models"
	"golang.org/x/crypto/bcrypt"
)

type MatchService struct {
	matchRepo MatchRepository
}

func NewMatchService(mr MatchRepository) *MatchService {
	return &MatchService{matchRepo: mr}
}

type AdminService struct {
	adminRepo AdminRepository
}

func NewAdminService(ar AdminRepository) *AdminService {
	return &AdminService{adminRepo: ar}
}

func (s *AdminService) AuthenticateAdmin(email, password string) (*models.Admin, error) {
	admin, err := s.adminRepo.FindAdminByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return admin, nil
}

func (s *AdminService) GetDashboardStats() (map[string]interface{}, error) {
	// Здесь логика получения статистики
	stats := map[string]interface{}{
		"total_users":    s.adminRepo.CountMatches(),
		"total_matches":  s.adminRepo.CountMatches(),
		"active_matches": s.adminRepo.CountActiveMatches(),
	}
	return stats, nil
}

func (s *MatchService) CreateMatch(match *models.Match) (*models.Match, error) {
	if err := validateMatch(match); err != nil {
		return nil, err
	}
	return s.matchRepo.CreateMatch(match)
}

func (s *MatchService) GetAllMatches() ([]models.Match, error) {
	return s.matchRepo.GetAll()
}

func (s *MatchService) DeleteMatch(id uint) error {
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return err
	}

	if match.Status != models.MatchUpcoming {
		return errors.New("cannot delete started or completed match")
	}

	return s.matchRepo.Delete(id)
}

type BetService struct {
	betRepo    BetRepository
	walletRepo wallet.WalletRepository
	matchRepo  MatchRepository
}

func NewBetService(
	wr *wallet.WalletRepository,
	mr *MatchRepository,
	br *BetRepository,
) *BetService {
	return &BetService{
		walletRepo: *wr,
		matchRepo:  *mr,
		betRepo:    *br,
	}
}

func (s *BetService) SettleBets(matchID uint, winner string) (*SettleResult, error) {
	bets, err := s.betRepo.GetByMatchID(matchID)
	if err != nil {
		return nil, err
	}

	result := &SettleResult{
		TotalBets: len(bets),
	}

	for _, b := range bets {
		if b.Team == winner {
			// Начисляем выигрыш
			_, tx, err := s.walletRepo.Deposit(b.WalletID, b.PotentialWin)
			if err != nil {
				return nil, err
			}

			// Обновляем статус ставки
			err = s.matchRepo.UpdateStatus(b.ID, models.MatchStatus(models.BetWon))
			if err != nil {
				// Откатываем транзакцию при ошибке (реализуйте при необходимости)
				return nil, err
			}

			result.WinningBets++
			result.Transactions = append(result.Transactions, tx)
		} else {
			err = s.matchRepo.UpdateStatus(b.ID, models.MatchStatus(models.BetLost))
			if err != nil {
				return nil, err
			}
		}
	}

	// Обновляем статус матча
	err = s.matchRepo.UpdateStatus(matchID, models.MatchCompleted)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type SettleResult struct {
	TotalBets    int
	WinningBets  int
	Transactions []*models.Transaction
}

func validateMatch(match *models.Match) error {
	if match.Team1 == "" || match.Team2 == "" {
		return errors.New("both teams are required")
	}

	if match.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}

	return nil
}
