package apis

import (
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListSales
// @Summary List sales
// @Tags Sale
// @Produce json
// @Param json query SaleListRequest true "query"
// @Success 200 {object} SaleListResponse
// @Router /sales [get]
func ListSales(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var query SaleListRequest
	if err := ValidateQuery(c, &query); err != nil {
		return err
	}

	// construct querySet
	querySet := query.QuerySet(DB).Order(ToOrderString(query.OrderBy, query.Sort))
	if query.BookID != nil {
		querySet = querySet.Where("book_id = ?", *query.BookID)
	}
	if query.UserID != nil {
		querySet = querySet.Where("user_id = ?", *query.UserID)
	}
	if query.StartTime != nil {
		querySet = querySet.Where("created_at >= ?", *query.StartTime)
	}
	if query.EndTime != nil {
		querySet = querySet.Where("created_at <= ?", *query.EndTime)
	}

	querySet = querySet.Session(&gorm.Session{}) // mark as safe to reuse

	var sales []Sale
	if err := querySet.Find(&sales).Error; err != nil {
		return err
	}

	var pageTotal int64
	if err := querySet.Offset(-1).Limit(-1).Count(&pageTotal).Error; err != nil {
		return err
	}

	var response SaleListResponse
	if err := copier.Copy(&response.Sales, &sales); err != nil {
		return err
	}
	response.PageTotal = int(pageTotal)

	return c.JSON(response)
}

// GetASale
// @Summary Get a sale by id
// @Tags Sale
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} SaleResponse
// @Router /sales/{id} [get]
func GetASale(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var sale Sale
	if err := DB.First(&sale, c.Params("id")).Error; err != nil {
		return err
	}

	var saleResponse SaleResponse
	if err := copier.Copy(&saleResponse, &sale); err != nil {
		return err
	}

	return c.JSON(saleResponse)
}

// CreateASale
// @Summary Create a sale
// @Tags Sale
// @Accept json
// @Produce json
// @Param json body SaleCreateRequest true "body"
// @Success 201 {object} SaleResponse
// @Router /sales [post]
func CreateASale(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var body SaleCreateRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	var sale Sale
	if err := copier.Copy(&sale, &body); err != nil {
		return err
	}
	sale.UserID = user.ID

	if err := DB.Create(&sale).Error; err != nil {
		return err
	}

	var saleResponse SaleResponse
	if err := copier.Copy(&saleResponse, &sale); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(saleResponse)
}
