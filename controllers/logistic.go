package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLogistic godoc
// @Summary Create a logistic
// @Description Mencatat bantuan logistik
// @Tags Logistic
// @Accept json
// @Produce json
// @Param input body structs.LogisticInput true "Data bantuan logistik"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /logistics [post]
func CreateLogistic(c *gin.Context) {
	var input structs.LogisticInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	logistic := structs.Logistic{
		Type:       input.Type,
		Quantity:   input.Quantity,
		Status:     input.Status,
		DisasterID: input.DisasterID,
	}

	err := repository.CreateLogistic(database.DbConnection, logistic)
	if err != nil {
		if err.Error() == "invalid logistics status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status logistik tidak valid, hanya bisa 'available', 'distributed', atau 'out_of_stock'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat bantuan logistik",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bantuan logistik berhasil dicatat",
		"result": logistic,
	})
}

// GetAllLogistics godoc
// @Summary Get all logistics
// @Description Mendapatkan daftar bantuan logistik
// @Tags Logistic
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /logistics [get]
func GetAllLogistics(c *gin.Context) {
	logistics, err := repository.GetAllLogistics(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar bantuan logistik",
		})
		return
	}

	if len(logistics) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar bantuan logistik yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": logistics,
	})
}

// GetLogisticByID godoc
// @Summary Get logistic by ID
// @Description Mendapatkan bantuan logistik berdasarkan ID
// @Tags Logistic
// @Accept json
// @Produce json
// @Param id path int true "Logistic ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /logistics/{id} [get]
func GetLogisticByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	logistic, err := repository.GetLogisticByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Bantuan logistik tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": logistic,
	})
}

// UpdateLogistic godoc
// @Summary Update logistic
// @Description Memperbarui bantuan logistik
// @Tags Logistic
// @Accept json
// @Produce json
// @Param id path int true "Logistic ID"
// @Param input body structs.LogisticInput true "Data bantuan logistik"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /logistics/{id} [put]
func UpdateLogistic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.Logistic
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}
	input.ID = id

	err = repository.UpdateLogistic(database.DbConnection, input)
	if err != nil {
		if err.Error() == "invalid logistics status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status logistik tidak valid, hanya bisa 'available', 'distributed', atau 'out_of_stock'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate bantuan logistik",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bantuan logistik berhasil diperbarui",
	})
}

// DeleteLogistic godoc
// @Summary Delete logistic
// @Description Menghapus bantuan logistik
// @Tags Logistic
// @Accept json
// @Produce json
// @Param id path int true "Logistic ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /logistics/{id} [delete]
func DeleteLogistic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteLogistic(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus bantuan logistik",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bantuan logistik berhasil dihapus",
	})
}
