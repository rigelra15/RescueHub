package controllers

import (
	"RescueHub/database"
	"RescueHub/middlewares"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login user
// @Description Autentikasi user untuk mendapatkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body structs.Login true "User login credentials"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /users/login [post]
func Login(c *gin.Context) {
	var loginRequest structs.Login

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	user, err := repository.GetUserByEmail(database.DbConnection, loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email atau password salah",
		})
		return
	}

	if !repository.CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email atau password salah",
		})
		return
	}

	token, err := middlewares.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Mendapatkan semua user
// @Tags Users
// @Produce json
// @Success 200 {object} []structs.User
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": users,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Mendapatkan user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.User
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	user, err := repository.GetUserByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

// CreateUser godoc
// @Summary Create new user
// @Description Membuat user baru
// @Tags Users
// @Accept json
// @Produce json
// @Param user body structs.UserInput true "User object"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input structs.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	hashedPassword, err := repository.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal meng-hash password",
		})
		return
	}

	user := structs.User{
		Name:       input.Name,
		Email:   		input.Email,
		Password:   hashedPassword,
		Role:       input.Role,
		Contact:    input.Contact,
	}

	err = repository.CreateUser(database.DbConnection, user)
	if err != nil {
		if err.Error() == "username sudah digunakan" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err.Error() == "invalid user role" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Role pengguna tidak valid, hanya bisa 'admin', 'donor', atau 'user'",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat user",
		})
		return
	}

	

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Mengubah data user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body structs.UserInput true "User object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var input structs.UserInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	hashedPassword, err := repository.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal meng-hash password",
		})
		return
	}

	user := structs.User{
		ID:         id,
		Name: 		 	input.Name,
		Email:      input.Email,
		Password:   hashedPassword,
		Role:       input.Role,
		Contact:    input.Contact,
	}

	err = repository.UpdateUser(database.DbConnection, user)
	if err != nil {
		if err.Error() == "invalid user role" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Role pengguna tidak valid, hanya bisa 'admin', 'donor', atau 'user'",
			})
			return
		}

		if err.Error() == "username sudah digunakan" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err.Error() == "user tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memperbarui user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diperbarui",
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteUser(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil dihapus",
	})
}

// UpdateUserInfoWithoutEmail godoc
// @Summary Update user without email
// @Description Mengubah data user tanpa email
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body structs.UpdateUserInfoWithoutEmail true "User object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/info/{id} [put]
func UpdateUserInfoWithoutEmail(c *gin.Context) {
	var input structs.UpdateUserInfoWithoutEmail
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	user := structs.User{
		ID:         id,
		Name: 		 	input.Name,
		Role:       input.Role,
		Contact:    input.Contact,
	}

	err = repository.UpdateUserInfoWithoutEmail(database.DbConnection, user)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memperbarui user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diperbarui",
	})
	
}