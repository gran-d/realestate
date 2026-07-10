package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)


func JWTProtected(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).
                JSON(fiber.Map{"error": "токен отсутствует"})
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenStr == authHeader {
            return c.Status(fiber.StatusUnauthorized).
                JSON(fiber.Map{"error": "формат: Bearer <токен>"})
        }

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).
                JSON(fiber.Map{"error": "токен невалиден или истёк"})
        }

        claims := token.Claims.(jwt.MapClaims)
        userID := int(claims["user_id"].(float64))
        c.Locals("user_id", userID)

        return c.Next()
    }
}