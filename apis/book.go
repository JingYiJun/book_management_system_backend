package apis

import (
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListBooks godoc
// @Summary List books
// @Tags Book
// @Accept json
// @Produce json
// @Param json query BookListRequest true "query"
// @Success 200 {object} BookListResponse
// @Router /books [get]
func ListBooks(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var query BookListRequest
	if err := ValidateQuery(c, &query); err != nil {
		return err
	}

	querySet := query.QuerySet(DB).Order(ToOrderString(query.OrderBy, query.Sort))
	if query.ID != nil {
		querySet = querySet.Where("id = ?", *query.ID)
	} else if query.ISBN != nil {
		querySet = querySet.Where("isbn = ?", *query.ISBN)
	} else {
		if query.Title != nil {
			querySet = querySet.Where("title LIKE ?", "%"+*query.Title+"%")
		}
		if query.Author != nil {
			querySet = querySet.Where("author LIKE ?", "%"+*query.Author+"%")
		}
		if query.Press != nil {
			querySet = querySet.Where("press LIKE ?", "%"+*query.Press+"%")
		}
		if query.OnSale != nil {
			querySet = querySet.Where("on_sale = ?", *query.OnSale)
		}
	}

	querySet = querySet.Session(&gorm.Session{}) // mark as safe to reuse

	var books []Book
	if err := querySet.Find(&books).Error; err != nil {
		return err
	}

	var pageTotal int64
	if err := querySet.Model(&Book{}).Offset(-1).Limit(-1).Count(&pageTotal).Error; err != nil {
		return err
	}

	var response BookListResponse
	if err := copier.Copy(&response.Books, &books); err != nil {
		return err
	}
	response.PageTotal = int(pageTotal)

	return c.JSON(response)
}

// GetABook godoc
// @Summary Get a book by id/isbn
// @Tags Book
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} BookResponse
// @Router /books/{id} [get]
func GetABook(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var comparedKeys = []string{"id", "isbn"}

	value := c.Params("id")
	if value == "" {
		return BadRequest()
	}

	var book Book
	for _, key := range comparedKeys {
		err := DB.Where("? = ?", key, value).First(&book).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			break
		}
	}
	if book.ID == 0 {
		return NotFound()
	}

	var bookResponse BookResponse
	if err := copier.Copy(&bookResponse, &book); err != nil {
		return err
	}

	return c.JSON(&bookResponse)
}

// CreateABook godoc
// @Summary Create a book
// @Tags Book
// @Accept json
// @Produce json
// @Param json body BookCreateRequest true "body"
// @Success 201 {object} BookResponse
// @Router /books [post]
func CreateABook(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	var body BookCreateRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	var book Book
	if err := copier.CopyWithOption(&book, &body, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	book.UserID = user.ID
	if err := DB.Create(&book).Error; err != nil {
		return err
	}

	var bookResponse BookResponse
	if err := copier.Copy(&bookResponse, &book); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(&bookResponse)
}

// ModifyABook godoc
// @Summary Modify a book
// @Tags Book
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param json body BookModifyRequest true "body"
// @Success 200 {object} BookResponse
// @Router /books/{id} [patch]
func ModifyABook(c *fiber.Ctx) error {
	var user User
	if err := GetCurrentUser(c, &user); err != nil {
		return err
	}

	bookID, err := c.ParamsInt("id")
	if err != nil {
		return BadRequest()
	}

	var book Book
	if err := DB.Where("id = ?", bookID).First(&book).Error; err != nil {
		return err
	}

	var body BookModifyRequest
	if err := ValidateBody(c, &body); err != nil {
		return err
	}

	if err := copier.CopyWithOption(&book, &body, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	if err := DB.Save(&book).Error; err != nil {
		return err
	}

	var bookResponse BookResponse
	if err := copier.Copy(&bookResponse, &book); err != nil {
		return err
	}

	return c.JSON(&bookResponse)
}
