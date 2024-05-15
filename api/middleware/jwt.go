package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT Auth --")

	token := c.Get("X-Api-Token")
	if token == "" {
		return fmt.Errorf("Unauthorized")
	}

	claims, err := parseJWT(token)
	if err != nil {
		return err
	}

	expiresAtFloat, ok := claims["expiresAt"].(float64)
	if !ok {
		fmt.Println("wrong format")
		return fmt.Errorf("Unauthorized")
	}
	expiresAt := int64(expiresAtFloat)

	expired := time.Now().Unix() > expiresAt
	if expired {
		return fmt.Errorf("Token expired")
	}

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
