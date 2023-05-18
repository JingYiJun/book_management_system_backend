package apis

import (
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListPurchases godoc
// @Summary List purchases
// @Tags Purchase
// @Produce json
// @Param json query PurchaseListRequest true "query"
// @Success 200 {object} PurchaseListResponse
// @Router /purchases [get]
func ListPurchases(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var query PurchaseListRequest
	if err := ValidateQuery(c, &query); err != nil {
		return err
	}

	querySet := query.QuerySet(DB).Order(ToOrderString(query.OrderBy, query.Sort))
	if query.BookID != nil {
		querySet = querySet.Where("book_id = ?", *query.BookID)
	}
	if query.UserID != nil {
		querySet = querySet.Where("user_id = ?", *query.UserID)
	}

	querySet = querySet.Session(&gorm.Session{}) // mark as safe to reuse

	var purchases []Purchase
	if err := querySet.Preload("Book").Find(&purchases).Error; err != nil {
		return err
	}

	var pageTotal int64
	if err := querySet.Model(&Purchase{}).Offset(-1).Limit(-1).Count(&pageTotal).Error; err != nil {
		return err
	}

	var response PurchaseListResponse
	if err := copier.Copy(&response.Purchases, &purchases); err != nil {
		return err
	}
	for i := range response.Purchases {
		if err := copier.Copy(&response.Purchases[i].Book, &purchases[i].Book); err != nil {
			return err
		}
	}
	response.PageTotal = int(pageTotal)

	return c.JSON(response)
}

// GetAPurchase godoc
// @Summary Get a purchase by id
// @Tags Purchase
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} PurchaseResponse
// @Router /purchases/{id} [get]
func GetAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var purchase Purchase
	if err := DB.First(&purchase, c.Params("id")).Error; err != nil {
		return NotFound()
	}

	var purchaseResponse PurchaseResponse
	if err := copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.JSON(&purchaseResponse)
}

// CreateAPurchase godoc
// @Summary Create a purchase
// @Tags Purchase
// @Accept json
// @Produce json
// @Param json body PurchaseCreateRequest true "body"
// @Success 201 {object} PurchaseResponse
// @Router /purchases [post]
func CreateAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var body PurchaseCreateRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	var purchase Purchase
	if err := copier.Copy(&purchase, &body); err != nil {
		return err
	}
	purchase.UserID = user.ID

	if err := DB.Create(&purchase).Error; err != nil {
		return err
	}

	var purchaseResponse PurchaseResponse
	if err := copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&purchaseResponse)
}

// ModifyAPurchase godoc
// @Summary Modify a purchase
// @Description Modify the quantity or price of a purchase by id
// @Tags Purchase
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param json body PurchaseModifyRequest true "body"
// @Success 200 {object} PurchaseResponse
// @Router /purchases/{id} [patch]
func ModifyAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	purchaseID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body PurchaseModifyRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	var purchase Purchase
	if err := DB.First(&purchase, purchaseID).Error; err != nil {
		return err
	}

	if purchase.Paid {
		return BadRequest("Cannot modify a paid purchase")
	}
	if purchase.Returned {
		return BadRequest("Cannot modify a returned purchase")
	}
	if purchase.Arrived {
		return BadRequest("Cannot modify an arrived purchase")
	}

	if err := copier.CopyWithOption(&purchase, &body, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	if err := DB.Save(&purchase).Error; err != nil {
		return err
	}

	var purchaseResponse PurchaseResponse
	if err := copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.JSON(&purchaseResponse)
}

// PayAPurchase godoc
// @Summary Pay a purchase
// @Description Pay a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} PurchaseResponse
// @Router /purchases/{id}/_pay [post]
func PayAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	purchaseID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var purchase Purchase
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Clauses(LockClause).First(&purchase, purchaseID).Error; err != nil {
			return err
		}

		if purchase.Paid {
			return BadRequest("Purchase has been paid")
		}
		if purchase.Returned {
			return BadRequest("Purchase has been returned")
		}

		purchase.Paid = true
		if err = tx.Model(&purchase).Update("paid", true).Error; err != nil {
			return err
		}

		balance := Balance{
			UserID:        user.ID,
			Change:        -purchase.Price * purchase.Quantity,
			OperationType: OperationTypePurchase,
			OperationID:   purchase.ID,
		}
		return tx.Create(&balance).Error
	})
	if err != nil {
		return err
	}

	var purchaseResponse PurchaseResponse
	if err := copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.JSON(&purchaseResponse)
}

// ReturnAPurchase godoc
// @Summary Return a purchase
// @Description Return a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} PurchaseResponse
// @Router /purchases/{id}/_return [post]
func ReturnAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	purchaseID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var purchase Purchase
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Clauses(LockClause).First(&purchase, purchaseID).Error; err != nil {
			return err
		}
		if purchase.Paid {
			return BadRequest("Purchase has been paid")
		}

		purchase.Returned = true
		return tx.Model(&purchase).Update("returned", true).Error
	})
	if err != nil {
		return err
	}

	var purchaseResponse PurchaseResponse
	if err = copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.JSON(&purchaseResponse)
}

// ArriveAPurchase
// @Summary Arrive a purchase
// @Description Arrive a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} PurchaseResponse
// @Router /purchases/{id}/_arrive [post]
func ArriveAPurchase(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	purchaseID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var purchase Purchase
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Clauses(LockClause).Preload("Book").First(&purchase, purchaseID).Error; err != nil {
			return err
		}

		if !purchase.Paid {
			return BadRequest("Purchase has not been paid")
		}
		if purchase.Returned {
			return BadRequest("Purchase has been returned")
		}

		purchase.Arrived = true
		if err = tx.Model(&purchase).Update("arrived", true).Error; err != nil {
			return err
		}

		// update book stock
		return tx.Model(&purchase.Book).Update("stock", gorm.Expr("stock + ?", purchase.Quantity)).Error
	})
	if err != nil {
		return err
	}

	var purchaseResponse PurchaseResponse
	if err := copier.Copy(&purchaseResponse, &purchase); err != nil {
		return err
	}

	return c.JSON(&purchase)
}
