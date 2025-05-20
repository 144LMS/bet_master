package admin

import (
	"github.com/144LMS/bet_master/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MatchRepository) CreateMatch(match *models.Match) (*models.Match, error) {
	if err := r.db.Create(match).Error; err != nil {
		return nil, err
	}
	return match, nil
}

func (r *MatchRepository) GetAll() ([]models.Match, error) {
	var matches []models.Match
	if err := r.db.Order("start_time asc").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *MatchRepository) GetByID(id uint) (*models.Match, error) {
	var match models.Match
	if err := r.db.First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}
func (r *AdminRepository) CountMatches() int64 {
	var count int64
	_ = r.db.Model(&models.Match{}).Count(&count).Error
	return count
}

func (r *AdminRepository) CountActiveMatches() int64 {
	var count int64
	_ = r.db.Model(&models.Match{}).Where("status IN ?", []string{"upcoming", "live"}).Count(&count).Error
	return count
}

func (r *AdminRepository) CountUsers() int64 {
	var count int64
	_ = r.db.Model(&models.User{}).Count(&count).Error
	return count
}

func (r *MatchRepository) Delete(id uint) error {
	return r.db.Delete(&models.Match{}, id).Error
}

func (r *MatchRepository) UpdateStatus(id uint, status models.MatchStatus) error {
	return r.db.Model(&models.Match{}).
		Where("id = ?", id).
		Update("status", status).Error
}

type BetRepository struct {
	db *gorm.DB
}

func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{db: db}
}

func (r *BetRepository) GetByMatchID(matchID uint) ([]models.Bet, error) {
	var bets []models.Bet
	err := r.db.Where("match_id = ?", matchID).Find(&bets).Error
	return bets, err
}

func (r *BetRepository) UpdateStatus(id uint, status models.BetStatus) error {
	return r.db.Model(&models.Bet{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *AdminRepository) CreateAdmin(admin *models.Admin) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)
	return r.db.Create(admin).Error
}
