package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/144LMS/bet_master/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AdminController struct {
	matchServ  MatchService
	betService BetService
	adminServ  AdminService
	jwtSecret  string
}

func NewAdminController(
	bs *BetService,
	ms *MatchService,
	as *AdminService,
	jwtSecret string,
) *AdminController {
	return &AdminController{
		betService: *bs,
		matchServ:  *ms,
		adminServ:  *as,
		jwtSecret:  jwtSecret,
	}
}

func (c *AdminController) AdminLogin(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := c.adminServ.AuthenticateAdmin(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  admin.ID,
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.SetCookie(
		"access_token",
		tokenString,
		3600*24, // срок действия: 1 день
		"/",
		"",    // domain
		false, // secure (поставь true если https)
		true,  // httpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user":    admin,
	})
	/*ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    admin.ID,
			"email": admin.Email,
			"role":  "admin",
		},
	})*/
}

func (c *AdminController) AdminDashboard(ctx *gin.Context) {
	// Получаем данные для дашборда (например, статистику)
	stats, err := c.adminServ.GetDashboardStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"stats":   stats,
		"message": "Welcome to admin dashboard",
	})
}

func (c *AdminController) CreateMatch(ctx *gin.Context) {
	var request struct {
		Team1     string    `json:"team1" binding:"required"`
		Team2     string    `json:"team2" binding:"required"`
		StartTime time.Time `json:"start_time" binding:"required"`
		SportType string    `json:"sport_type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match := models.Match{
		Team1:     request.Team1,
		Team2:     request.Team2,
		StartTime: request.StartTime,
		SportType: request.SportType,
		Status:    models.MatchUpcoming,
	}

	createdMatch, err := c.matchServ.CreateMatch(&match)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Match created successfully",
		"match":   createdMatch,
	})
}

func (c *AdminController) GetAllMatches(ctx *gin.Context) {
	matches, err := c.matchServ.GetAllMatches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"matches": matches, // Теперь поля будут правильно сериализованы
	})
}

func (c *AdminController) DeleteMatch(ctx *gin.Context) {
	matchID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err := c.matchServ.DeleteMatch(uint(matchID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "match deleted"})
}

func (c *AdminController) SettleMatch(ctx *gin.Context) {
	matchID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var request SettleMatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.betService.SettleBets(uint(matchID), request.Winner)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "bets settled",
		"total_bets":   result.TotalBets,
		"winning_bets": result.WinningBets,
	})
}

func (c *AdminController) CreateAdmin(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin := models.Admin{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := c.adminServ.adminRepo.CreateAdmin(&admin); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully"})
}

// Вспомогательные структуры запросов/ответов
type CreateMatchRequest struct {
	Team1     string    `json:"team1" binding:"required"`
	Team2     string    `json:"team2" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	SportType string    `json:"sport_type" binding:"required"`
}

type SettleMatchRequest struct {
	Winner string `json:"winner" binding:"required,oneof=team1 team2 draw"`
}

type SettleMatchResponse struct {
	TotalBets   int `json:"total_bets"`
	WinningBets int `json:"winning_bets"`
}
