package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"
)

// SecretKey todo: remove SecretKey from here
const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.RegisteredClaims
	Scope string
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "user is unauthorized",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)

	if !strings.Contains(c.Path(), payload.Scope) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "invalid scope",
		})
	}

	return c.Next()
}

func GetUserId(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*ClaimsWithScope)

	userId, err := strconv.Atoi(payload.Subject)

	return uint(userId), nil
}

func GenerateJWT(userId uint, scope string) (string, error) {
	claims := ClaimsWithScope{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userId)),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		},
		Scope: scope,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
}
