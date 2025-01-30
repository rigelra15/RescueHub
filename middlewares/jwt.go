package middlewares

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"log"
	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateSecretKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalf("Error generating secret key: %v", err)
	}
	return key
}

var secretKey = generateSecretKey()

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(secretKey)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("authHeader: ", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Gagal memproses token: %v", err.Error()),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Gagal memproses klaim token",
			})
			c.Abort()
			return
		}

		email, exists := claims["email"].(string)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak mengandung email",
			})
			c.Abort()
			return
		}

		c.Set("email", email)

		c.Next()
	}
}

