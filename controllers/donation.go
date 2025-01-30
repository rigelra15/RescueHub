package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDonation godoc
// @Summary Create a donation
// @Description Mencatat donasi
// @Tags Donation
// @Accept json
// @Produce json
// @Param input body structs.DonationInput true "Data donasi"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /donations [post]
func CreateDonation(c *gin.Context) {
	var input structs.DonationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	donation := structs.Donation{
		DonorID:    input.DonorID,
		DisasterID: input.DisasterID,
		Amount:     input.Amount,
		ItemName:   input.ItemName,
		Status:     input.Status,
	}

	err := repository.CreateDonation(database.DbConnection, donation)
	if err != nil {
		if err.Error() == "invalid donation status" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status donasi tidak valid, hanya bisa 'pending', 'confirmed', atau 'rejected'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat donasi",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Donasi berhasil dicatat",
	})
}

// GetAllDonations godoc
// @Summary Get all donations
// @Description Mendapatkan daftar donasi
// @Tags Donation
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /donations [get]
func GetAllDonations(c *gin.Context) {
	donations, err := repository.GetAllDonations(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar donasi yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": donations,
	})
}

// GetDonationByID godoc
// @Summary Get a donation by ID
// @Description Mendapatkan donasi berdasarkan ID
// @Tags Donation
// @Accept json
// @Produce json
// @Param id path int true "Donation ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /donations/{id} [get]
func GetDonationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	donation, err := repository.GetDonationByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Donasi tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": donation,
	})
}

// UpdateDonation godoc
// @Summary Update a donation
// @Description Memperbarui donasi
// @Tags Donation
// @Accept json
// @Produce json
// @Param id path int true "Donation ID"
// @Param input body structs.DonationInput true "Data donasi"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /donations/{id} [put]
func UpdateDonation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.DonationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}
	
	donation := structs.Donation{
		ID:         id,
		DonorID:    input.DonorID,
		DisasterID: input.DisasterID,
		Amount:     input.Amount,
		ItemName:   input.ItemName,
		Status:     input.Status,
	}

	err = repository.UpdateDonation(database.DbConnection, donation)
	if err != nil {
		if err.Error() == "invalid donation status" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status donasi tidak valid, hanya bisa 'pending', 'confirmed', atau 'rejected'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate donasi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Donasi berhasil diperbarui",
	})
}

// DeleteDonation godoc
// @Summary Delete a donation
// @Description Menghapus donasi
// @Tags Donation
// @Accept json
// @Produce json
// @Param id path int true "Donation ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /donations/{id} [delete]
func DeleteDonation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteDonation(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus donasi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Donasi berhasil dihapus",
	})
}
