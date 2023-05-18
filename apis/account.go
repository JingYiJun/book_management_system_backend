package apis

import (
	"book_management_system_backend/config"
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

var (
	ErrInvalidUsernameOrPassword = BadRequest("用户名或密码错误")
	ErrUserAlreadyExist          = BadRequest("用户名已存在")
)

// Register godoc
// @Summary Register, admin only
// @Tags Account
// @Accept json
// @Produce json
// @Param json body RegisterRequest true "body"
// @Success 201 {object} User
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	var currentUser User
	err := GetCurrentUser(c, &currentUser)
	if err != nil {
		return err
	}

	if !currentUser.IsAdmin {
		return Forbidden("Only admin can register new user")
	}

	var body RegisterRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var user User
	err = copier.CopyWithOption(&user, &body, copier.Option{IgnoreEmpty: true})
	if err != nil {
		return err
	}
	user.HashedPassword = MakePassword(body.Password)

	result := DB.Where(User{Username: body.Username}).Attrs(user).FirstOrCreate(&user)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrUserAlreadyExist
	}

	token, err := GenerateToken(&user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "access",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/api",
		Domain:  config.Config.Hostname,
	})

	return c.JSON(TokenResponse{
		AccessToken: token,
		Message:     "注册成功",
	})
}

// Login godoc
// @Summary Login
// @Tags Account
// @Accept json
// @Produce json
// @Param json body LoginRequest true "body"
// @Success 200 {object} User
// @Router /login [post]
func Login(c *fiber.Ctx) error {
	var body LoginRequest
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var user User
	err = DB.Take(&user, "username = ?", body.Username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidUsernameOrPassword
		} else {
			return err
		}
	}

	if !CheckPassword(body.Password, user.HashedPassword) {
		return ErrInvalidUsernameOrPassword
	}

	token, err := GenerateToken(&user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "access",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/api",
		Domain:  config.Config.Hostname,
	})

	return c.JSON(TokenResponse{
		AccessToken: token,
		Message:     "登录成功",
	})
}

// GetUserMe godoc
// @Summary Get current user
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /users/me [get]
func GetUserMe(c *fiber.Ctx) error {
	var user User
	err := GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	err = DB.Take(&user).Error
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// ModifyUserMe godoc
// @Summary modify current user
// @Tags Account
// @Accept json
// @Produce json
// @Param json body UserModifyRequest true "body"
// @Success 200 {object} User
// @Router /users/me [patch]
func ModifyUserMe(c *fiber.Ctx) error {
	var user User
	err := GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	var body UserModifyRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(LockClause).Take(&user, user.ID).Error
		if err != nil {
			return err
		}

		err = copier.CopyWithOption(&user, &body, copier.Option{IgnoreEmpty: true})
		if err != nil {
			return err
		}

		if body.Password != nil {
			user.HashedPassword = MakePassword(*body.Password)
		}

		return tx.Save(&user).Error
	})
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// DeleteUserMe godoc
// @Summary delete self
// @Tags Account
// @Accept json
// @Produce json
// @Success 204
// @Router /users/me [delete]
func DeleteUserMe(c *fiber.Ctx) error {
	var user User
	err := GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	if user.ID == 1 {
		return Forbidden("Can't delete first admin")
	}

	err = DB.Delete(&user).Error
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers godoc
// @Summary list users, admin only
// @Tags Account
// @Accept json
// @Produce json
// @Param page query UserListRequest true "page"
// @Success 200 {object} UserListResponse
// @Router /users [get]
func ListUsers(c *fiber.Ctx) error {
	var currentUser User
	err := GetCurrentUser(c, &currentUser)
	if err != nil {
		return err
	}
	if !currentUser.IsAdmin {
		return Forbidden()
	}

	var query UserListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	querySet := query.QuerySet(DB).Order(query.OrderBy + " " + query.Sort).Session(&gorm.Session{}) // mark as safe to reuse

	var users []User
	err = querySet.Find(&users).Error
	if err != nil {
		return err
	}

	var pageTotal int64
	err = querySet.Model(&User{}).Limit(-1).Offset(-1).Count(&pageTotal).Error
	if err != nil {
		return err
	}

	var response UserListResponse
	if err := copier.Copy(&response.Users, &users); err != nil {
		return err
	}
	response.PageTotal = int(pageTotal)

	return c.JSON(response)
}

// GetUser godoc
// @Summary get a user by id/username/staff_id, admin only
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	var currentUser User
	if err := GetCurrentUser(c, &currentUser); err != nil {
		return err
	}
	if !currentUser.IsAdmin {
		return Forbidden()
	}

	value := c.Params("id")
	if value == "" {
		return BadRequest()
	}

	var comparedKeys = []string{"id", "username", "staff_id"}
	var user User
	for _, key := range comparedKeys {
		err := DB.Where("? = ?", key, value).First(&user).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			break
		}
	}

	return c.JSON(user)
}

// ModifyAUser godoc
// @Summary modify a user by id, admin only
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} User
// @Router /users/{id} [patch]
// @Param body body UserModifyRequest true "body"
func ModifyAUser(c *fiber.Ctx) error {
	var currentUser User
	err := GetCurrentUser(c, &currentUser)
	if err != nil {
		return err
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body UserModifyRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var user User
	err = DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(LockClause).Take(&user, userID).Error
		if err != nil {
			return err
		}

		if !(user.ID == currentUser.ID || currentUser.IsAdmin) {
			return Forbidden()
		}

		err = copier.CopyWithOption(&user, &body, copier.Option{IgnoreEmpty: true})
		if err != nil {
			return err
		}

		if body.Password != nil {
			user.HashedPassword = MakePassword(*body.Password)
		}

		return tx.Save(&user).Error
	})
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// DeleteAUser godoc
// @Summary delete a user by id, admin only
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204
// @Router /users/{id} [delete]
func DeleteAUser(c *fiber.Ctx) error {
	var currentUser User
	err := GetCurrentUser(c, &currentUser)
	if err != nil {
		return err
	}
	if !currentUser.IsAdmin {
		return Forbidden()
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if userID == 1 {
		return Forbidden("Can't delete first admin")
	}

	var user User
	err = DB.Delete(&user, userID).Error
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
