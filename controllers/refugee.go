package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRefugee godoc
// @Summary Create a refugee
// @Description Membuat data pengungsi baru
// @Tags Refugee
// @Accept json
// @Produce json
// @Param input body structs.RefugeeInput true "Data pengungsi"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /refugees [post]
func CreateRefugee(c *gin.Context) {
	var input structs.RefugeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	refugee := structs.Refugee{
		Name:       input.Name,
		Age:        input.Age,
		Condition:  input.Condition,
		Needs:      input.Needs,
		ShelterID:  input.ShelterID,
		DisasterID: input.DisasterID,
	}

	err := repository.CreateRefugee(database.DbConnection, refugee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat data pengungsi",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data pengungsi berhasil dibuat",
	})
}

// GetAllRefugees godoc
// @Summary Get all refugees
// @Description Mendapatkan daftar pengungsi
// @Tags Refugee
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /refugees [get]
func GetAllRefugees(c *gin.Context) {
	refugees, err := repository.GetAllRefugees(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar pengungsi",
		})
		return
	}

	if len(refugees) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar pengungsi yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": refugees,
	})
}

// GetRefugeeByID godoc
// @Summary Get refugee by ID
// @Description Mendapatkan data pengungsi berdasarkan ID
// @Tags Refugee
// @Accept json
// @Produce json
// @Param id path int true "Refugee ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /refugees/{id} [get]
func GetRefugeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	refugee, err := repository.GetRefugeeByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pengungsi tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": refugee,
	})
}

// UpdateRefugee godoc
// @Summary Update a refugee
// @Description Memperbarui data pengungsi
// @Tags Refugee
// @Accept json
// @Produce json
// @Param id path int true "Refugee ID"
// @Param input body structs.RefugeeInput true "Data pengungsi"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /refugees/{id} [put]
func UpdateRefugee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.RefugeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	refugee := structs.Refugee{
		ID:         id,
		Name:       input.Name,
		Age:        input.Age,
		Condition:  input.Condition,
		Needs:      input.Needs,
		ShelterID:  input.ShelterID,
		DisasterID: input.DisasterID,
	}

	err = repository.UpdateRefugee(database.DbConnection, refugee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate data pengungsi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pengungsi berhasil diperbarui",
	})
}

// DeleteRefugee godoc
// @Summary Delete a refugee
// @Description Menghapus data pengungsi
// @Tags Refugee
// @Accept json
// @Produce json
// @Param id path int true "Refugee ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /refugees/{id} [delete]
func DeleteRefugee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteRefugee(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus data pengungsi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pengungsi berhasil dihapus",
	})
}

// GetDistributionLogsByRefugeeID godoc
// @Summary Get distribution logs by refugee ID
// @Description Menampilkan daftar distribusi bantuan yang diterima oleh pengungsi tertentu
// @Tags Refugee
// @Produce json
// @Param id path int true "Refugee ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /refugees/{id}/distribution-logs [get]
func GetDistributionLogsByRefugeeID(c *gin.Context) {
	refugeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengungsi tidak valid"})
		return
	}

	logs, err := repository.GetDistributionLogsByRefugeeID(database.DbConnection, refugeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan catatan distribusi bantuan"})
		return
	}

	if len(logs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada distribusi bantuan untuk pengungsi ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": logs})
}