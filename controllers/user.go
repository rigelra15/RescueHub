package controllers

import (
	"RescueHub/database"
	"RescueHub/middlewares"
	"RescueHub/repository"
	"RescueHub/structs"
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	user, err := repository.GetUserByEmail(database.DbConnection, loginRequest.Email)
	if err != nil || !repository.CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if user.Role == "admin" || user.Is2FA {
		fmt.Println("User ID yang digunakan untuk OTP:", user.ID)
		otp := middlewares.GenerateOTP()
		err = middlewares.SendOTPToEmail(user.Email, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim OTP"})
			return
		}

		err = repository.SaveOTP(database.DbConnection, user.ID, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP telah dikirim ke email Anda"})
		return
	}

	token, err := middlewares.GenerateJWT(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verifikasi OTP untuk mendapatkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body structs.VerifyOTP true "User OTP verification"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /users/verify-otp [post]
func VerifyOTP(c *gin.Context) {
	var otpRequest struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&otpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	user, err := repository.GetUserByEmail(database.DbConnection, otpRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	isValid, err := repository.ValidateOTP(database.DbConnection, user.ID, otpRequest.OTP)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP salah atau telah kedaluwarsa"})
		return
	}

	token, err := middlewares.GenerateJWT(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Enable2FA godoc
// @Summary Enable 2FA
// @Description Mengaktifkan atau menonaktifkan 2FA
// @Tags Users
// @Accept json
// @Produce json
// @Param user body structs.Enable2FA true "User 2FA status"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/enable-2fa [put]
func Enable2FA(c *gin.Context) {
	var request struct {
		Is2FAEnabled bool `json:"is_2fa"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak dapat mengidentifikasi pengguna"})
		return
	}

	err := repository.Enable2FA(database.DbConnection, email.(string), request.Is2FAEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status 2FA"})
		return
	}

	status := "dinonaktifkan"
	if request.Is2FAEnabled {
		status = "diaktifkan"
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA berhasil " + status})
}

// ChangeUserRole godoc
// @Summary Mengubah role pengguna
// @Description Admin dapat mengubah role siapa saja, sedangkan user hanya bisa mengubah role sendiri ke donor atau user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param role body structs.ChangeUserRole true "List role: admin, donor, user"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 403 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id}/change-role [put]

func isValidUserRole(role string) bool {
	validRoles := []string{"admin", "donor", "user"}
	for _, valid := range validRoles {
		if role == valid {
			return true
		}
	}
	return false
}

func ChangeUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role"`
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if !isValidUserRole(input.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid, hanya bisa 'admin', 'donor', atau 'user'"})
		return
	}

	authEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak dapat mengidentifikasi pengguna"})
		return
	}

	currentUser, err := repository.GetUserByEmail(database.DbConnection, authEmail.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if currentUser.ID != id && currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengubah role pengguna lain"})
		return
	}

	if currentUser.Role != "admin" && input.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengubah role ke admin"})
		return
	}

	is2FA := input.Role == "admin"

	err = repository.UpdateUserRole(database.DbConnection, id, input.Role, is2FA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role pengguna berhasil diperbarui"})
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
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input structs.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	adminCount, err := repository.CountAdmins(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memeriksa jumlah admin",
		})
		return
	}

	if adminCount > 0 && input.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Pendaftaran dengan role admin tidak diizinkan",
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

	is2FA := false
	if input.Role == "admin" {
		is2FA = true
	}

	user := &structs.User{
		Name:    input.Name,
		Email:   input.Email,
		Password: hashedPassword,
		Role:    input.Role,
		Contact: input.Contact,
		Is2FA:   is2FA,
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
		"result":  user,
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

	authEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak dapat mengidentifikasi pengguna"})
		return
	}

	currentUser, err := repository.GetUserByEmail(database.DbConnection, authEmail.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if currentUser.ID == id && input.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak bisa mengubah role Anda sendiri menjadi admin"})
		return
	}

	if currentUser.Role != "admin" && input.Role != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengubah role"})
		return
	}

	var hashedPassword string
	if input.Password != "" {
		hashedPassword, err = repository.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal meng-hash password",
			})
			return
		}
	} else {
		hashedPassword = currentUser.Password
	}

	user := structs.User{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      input.Role,
		Contact:   input.Contact,
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
// @Router /users/{id}/edit-info [put]
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

// GetDonationsByUserID godoc
// @Summary Get donations by user ID
// @Description Menampilkan daftar donasi yang diberikan oleh user tertentu
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id}/donations [get]
func GetDonationsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	donations, err := repository.GetDonationsByUserID(database.DbConnection, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan daftar donasi"})
		return
	}

	if len(donations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada donasi yang tercatat untuk user ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": donations})
}

// GetEmergencyReportsByUserID godoc
// @Summary Get emergency reports by user ID
// @Description Menampilkan semua laporan darurat yang dikirim oleh pengguna tertentu
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /users/{id}/emergency-reports [get]
func GetEmergencyReportsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	reports, err := repository.GetEmergencyReportsByUserID(database.DbConnection, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan laporan darurat"})
		return
	}

	if len(reports) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada laporan darurat yang dibuat oleh pengguna ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": reports})
}