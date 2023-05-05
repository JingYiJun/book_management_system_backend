package apis

import (
	. "book_management_system_backend/models"
	. "book_management_system_backend/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	ID      int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
}

func GetCurrentUser(c *fiber.Ctx, user *User) error {
	accessToken := c.Cookies("access")
	if accessToken == "" {
		accessToken = c.Get("Authorization")
		if accessToken == "" {
			return Unauthorized()
		}
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}
	}
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if userClaims, ok := token.Claims.(*UserClaims); !ok {
			return nil, errors.New("invalid jwt token")
		} else {
			var userJwtSecret UserJwtSecret
			err := DB.Take(&userJwtSecret, userClaims.ID).Error
			if err != nil {
				return nil, err
			}
			return []byte(userJwtSecret.Secret), nil
		}
	})
	if err != nil {
		Logger.Error("invalid jwt token", zap.String("token", accessToken), zap.Error(err))
		return Unauthorized()
	}

	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user.ID = userClaims.ID
		user.IsAdmin = userClaims.IsAdmin
		c.Locals("user_id", user.ID)
		return nil
	} else {
		Logger.Error("invalid jwt token", zap.String("token", accessToken))
		return Unauthorized()
	}
}

func GenerateToken(user *User) (string, error) {
	var userJwtSecret UserJwtSecret
	err := DB.Take(&userJwtSecret, user.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userJwtSecret = UserJwtSecret{
				ID:     user.ID,
				Secret: randstr.Base62(32),
			}
			err = DB.Create(&userJwtSecret).Error
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID:      user.ID,
		IsAdmin: user.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(userJwtSecret.Secret))
}
