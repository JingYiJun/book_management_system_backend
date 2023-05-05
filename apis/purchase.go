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
// @Success 200 {array} Purchase
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

	var purchases []Purchase
	if err := querySet.Find(&purchases).Error; err != nil {
		return err
	}

	return c.JSON(purchases)
}

// GetAPurchase godoc
// @Summary Get a purchase by id
// @Tags Purchase
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} Purchase
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

	return c.JSON(&purchase)
}

// CreateAPurchase godoc
// @Summary Create a purchase
// @Tags Purchase
// @Accept json
// @Produce json
// @Param json body PurchaseCreateRequest true "body"
// @Success 201 {object} Purchase
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

	return c.JSON(&purchase)
}

// ModifyAPurchase godoc
// @Summary Modify a purchase
// @Description Modify the quantity or price of a purchase by id
// @Tags Purchase
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param json body PurchaseModifyRequest true "body"
// @Success 200 {object} Purchase
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

	return c.JSON(&purchase)
}

// PayAPurchase godoc
// @Summary Pay a purchase
// @Description Pay a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} Purchase
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
		if err := tx.Clauses(LockClause).First(&purchase, purchaseID).Error; err != nil {
			return err
		}

		if purchase.Paid {
			return BadRequest("Purchase has been paid")
		}
		if purchase.Returned {
			return BadRequest("Purchase has been returned")
		}

		if err := tx.Model(&purchase).Update("paid", true).Error; err != nil {
			return err
		}

		// expense

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(&purchase)
}

// ReturnAPurchase godoc
// @Summary Return a purchase
// @Description Return a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} Purchase
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
		if err := tx.Clauses(LockClause).First(&purchase, purchaseID).Error; err != nil {
			return err
		}

		if purchase.Paid {
			return BadRequest("Purchase has been paid")
		}

		if err := tx.Model(&purchase).Updates(fiber.Map{"returned": true, "arrived": false}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(&purchase)
}

// ArriveAPurchase
// @Summary Arrive a purchase
// @Description Arrive a purchase by id
// @Tags Purchase
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} Purchase
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
		if err := tx.Clauses(LockClause).Preload("Book").First(&purchase, purchaseID).Error; err != nil {
			return err
		}

		if !purchase.Paid {
			return BadRequest("Purchase has not been paid")
		}
		if purchase.Returned {
			return BadRequest("Purchase has been returned")
		}

		if err := tx.Model(&purchase).Update("arrived", true).Error; err != nil {
			return err
		}

		// update book stock
		if err := tx.Model(&purchase.Book).Update("stock", gorm.Expr("stock + ?", purchase.Quantity)).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(&purchase)
}
