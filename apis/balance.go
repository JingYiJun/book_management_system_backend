package apis

import (
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListBalances godoc
// @Summary List balances
// @Tags Balance
// @Produce json
// @Param json query BalanceListRequest true "query"
// @Success 200 {object} BalanceListResponse
// @Router /balances [get]
func ListBalances(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var query BalanceListRequest
	if err := ValidateQuery(c, &query); err != nil {
		return err
	}

	querySet := query.QuerySet(DB).Order(ToOrderString(query.OrderBy, query.Sort))
	if query.UserID != nil {
		querySet = querySet.Where("user_id = ?", *query.UserID)
	}
	if query.Positive != nil {
		if *query.Positive {
			querySet = querySet.Where("change > 0")
		} else {
			querySet = querySet.Where("change < 0")
		}
	}
	if query.StartTime != nil {
		querySet = querySet.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		querySet = querySet.Where("created_at <= ?", *query.EndTime)
	}

	querySet = querySet.Session(&gorm.Session{}) // mark as safe to reuse

	var balances []Balance
	if err := querySet.Find(&balances).Error; err != nil {
		return err
	}

	var pageTotal int64
	if err := querySet.Offset(-1).Limit(-1).Count(&pageTotal).Error; err != nil {
		return err
	}

	var response BalanceListResponse
	if err := copier.Copy(&response.Balances, &balances); err != nil {
		return err
	}
	response.PageTotal = int(pageTotal)

	return c.JSON(response)
}

// CreateABalance godoc
// @Summary Create a balance
// @Tags Balance
// @Accept json
// @Produce json
// @Param json body BalanceCreateRequest true "body"
// @Success 201 {object} BalanceResponse
// @Router /balances [post]
func CreateABalance(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var body BalanceCreateRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	balance := Balance{
		UserID:        user.ID,
		Change:        body.Change(),
		OperationType: OperationTypeManual,
		Reason:        body.Reason,
	}
	if err := DB.Create(&balance).Error; err != nil {
		return err
	}

	var balanceResponse BalanceResponse
	if err := copier.Copy(&balanceResponse, &balance); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&balanceResponse)
}

// GetABalance godoc
// @Summary Get a balance by id
// @Tags Balance
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} BalanceResponse
// @Router /balances/{id} [get]
func GetABalance(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var balance Balance
	if err := DB.First(&balance, c.Params("id")).Error; err != nil {
		return err
	}

	var balanceResponse BalanceResponse
	if err := copier.Copy(&balanceResponse, &balance); err != nil {
		return err
	}

	return c.JSON(&balanceResponse)
}
