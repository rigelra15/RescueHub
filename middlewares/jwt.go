package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"RescueHub/database"
	"RescueHub/repository"
	"strconv"
	"database/sql"
	"fmt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 3).Unix(),
	})

	return token.SignedString(secretKey)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
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

		email, emailExists := claims["email"].(string)
		role, roleExists := claims["role"].(string)

		if !emailExists || !roleExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak mengandung informasi pengguna",
			})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

func RequireSameUserOrRole(messages string, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, roleExists := c.Get("role")
		email, emailExists := c.Get("email")
		userIDParam := c.Param("id")
		userID, err := strconv.Atoi(userIDParam)

		if !roleExists || !emailExists || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak memiliki akses",
			})
			c.Abort()
			return
		}

		if role.(string) == requiredRole {
			c.Next()
			return
		}

		user, err := repository.GetUserByEmail(database.DbConnection, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Pengguna tidak ditemukan",
			})
			c.Abort()
			return
		}

		if user.ID != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": messages,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireSelfFor2FA() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak memiliki akses"})
			c.Abort()
			return
		}

		user, err := repository.GetUserByEmail(database.DbConnection, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
			c.Abort()
			return
		}

		if user.Role == "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin wajib menggunakan 2FA"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRoles(errorMessage string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses ditolak, role tidak ditemukan",
			})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses ditolak, format role tidak valid",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": errorMessage,
		})
		c.Abort()
	}
}

func RequireVolunteerOrRole(messages string, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, roleExists := c.Get("role")
		email, emailExists := c.Get("email")

		if !roleExists || !emailExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak memiliki akses",
			})
			c.Abort()
			return
		}

		if role.(string) == requiredRole {
			c.Next()
			return
		}

		user, err := repository.GetUserByEmail(database.DbConnection, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Pengguna tidak ditemukan",
			})
			c.Abort()
			return
		}

		isVolunteer, err := repository.IsUserVolunteer(database.DbConnection, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal memeriksa status volunteer",
			})
			c.Abort()
			return
		}

		if isVolunteer {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": messages,
		})
		c.Abort()
	}
}

func RequireSelfForRelatedEntities(messages, tableName, userColumn string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, roleExists := c.Get("role")
		email, emailExists := c.Get("email")
		if !emailExists || !roleExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak memiliki akses",
			})
			c.Abort()
			return
		}

		if role.(string) == "admin" {
			c.Next()
			return
		}

		user, err := repository.GetUserByEmail(database.DbConnection, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Pengguna tidak ditemukan",
			})
			c.Abort()
			return
		}

		entityIDParam := c.Param("id")
		entityID, err := strconv.Atoi(entityIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID tidak valid",
			})
			c.Abort()
			return
		}

		var ownerID int
		query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", userColumn, tableName)
		err = database.DbConnection.QueryRow(query, entityID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Data tidak ditemukan",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal mengambil data",
			})
			c.Abort()
			return
		}

		if user.ID != ownerID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": messages,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAdminOrSelfForRoleChange() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, roleExists := c.Get("role")
		email, emailExists := c.Get("email")
		userIDParam := c.Param("id")
		userID, err := strconv.Atoi(userIDParam)

		if !roleExists || !emailExists || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tidak memiliki akses",
			})
			c.Abort()
			return
		}

		currentUser, err := repository.GetUserByEmail(database.DbConnection, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
			c.Abort()
			return
		}

		if currentUser.ID != userID && currentUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengubah role pengguna lain"})
			c.Abort()
			return
		}

		c.Next()
	}
}
