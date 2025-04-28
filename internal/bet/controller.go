package bet

import (
	"net/http"

	"github.com/144LMS/bet_master/internal/wallet"
	"github.com/144LMS/bet_master/models"
	"github.com/gin-gonic/gin"
)

type BetController struct {
	betService    BetService
	walletService wallet.WalletService
}

func NewBetController(bs BetService, ws wallet.WalletService) *BetController {
	return &BetController{
		betService:    bs,
		walletService: ws,
	}
}

func (c *BetController) PlaceBet(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx) // ваша реализация получения ID

	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Odds   float64 `json:"odds" binding:"required,gt=1"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Списываем средства
	wallet, err := c.walletService.Withdraw(userID, request.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Недостаточно средств"})
		return
	}

	// 2. Создаём ставку
	bet, err := c.betService.PlaceBet(models.Bet{
		UserID:       userID,
		Amount:       request.Amount,
		Odds:         request.Odds,
		PotentialWin: request.Amount * request.Odds,
	})

	if err != nil {
		// Если ошибка - возвращаем деньги
		_, _ = c.walletService.Deposit(userID, request.Amount)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"bet":    bet,
		"wallet": wallet,
	})
}
