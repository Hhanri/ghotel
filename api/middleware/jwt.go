package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hhanri/ghotel/api"
	"github.com/hhanri/ghotel/db"
)

type JWTMiddleware struct {
	store *db.Store
}

func NewJWTMiddleware(store *db.Store) *JWTMiddleware {
	return &JWTMiddleware{
		store: store,
	}
}

func (m *JWTMiddleware) JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT Auth --")

	token := c.Get("X-Api-Token")
	if token == "" {
		return api.FiberUnauthorizedErrorResponse(c)
	}

	claims, err := parseJWT(token)
	if err != nil {
		return api.FiberUnauthorizedErrorResponse(c)
	}

	expiresAtFloat, ok := claims["expiresAt"].(float64)
	if !ok {
		fmt.Println("wrong format")
		return api.FiberUnauthorizedErrorResponse(c)
	}
	expiresAt := int64(expiresAtFloat)

	expired := time.Now().Unix() > expiresAt
	if expired {
		return api.FiberExpiredTokenErrorResponse(c)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return api.FiberUnauthorizedErrorResponse(c)
	}

	user, err := m.store.User.GetUserByID(c.Context(), id)
	if err != nil {
		return api.FiberUnauthorizedErrorResponse(c)
	}
	c.Context().SetUserValue("user", user)
	return c.Next()
}

func parseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT")
		return nil, fmt.Errorf("Unauthorized")
	}

	if !token.Valid {
		fmt.Println("Invalid Token")
		return nil, fmt.Errorf("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}

	return claims, nil
}
