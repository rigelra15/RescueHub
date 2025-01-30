package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllDisasters godoc
// @Summary Get all disasters
// @Description Mendapatkan semua laporan bencana
// @Tags Disaster
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /disasters [get]
func GetAllDisasters(c *gin.Context) {
	disasters, err := repository.GetAllDisasters(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan laporan bencana",
		})
		return
	}

	if len(disasters) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tidak ada laporan bencana",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": disasters,
	})
}


// GetDisasterByID godoc
// @Summary Get disaster by ID
// @Description Mendapatkan laporan bencana berdasarkan ID
// @Tags Disaster
// @Produce json
// @Param id path int true "Disaster ID"
// @Success 200 {object} structs.Disaster
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /disasters/{id} [get]
func GetDisasterByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	disaster, err := repository.GetDisasterByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Laporan bencana tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": disaster,
	})
}

// CreateDisaster godoc
// @Summary Create disaster
// @Description Membuat laporan bencana baru
// @Tags Disaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param disaster body structs.DisasterInput true "Disaster object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /disasters [post]
func CreateDisaster(c *gin.Context) {
	var input structs.DisasterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	disaster := structs.Disaster{
		Type:        input.Type,
		Location:    input.Location,
		Description: input.Description,
		Status:      input.Status,
		ReportedBy:  input.ReportedBy,
	}

	err := repository.CreateDisaster(database.DbConnection, disaster)
	if err != nil {
		if err.Error() == "disaster already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Nama bencana sudah terdaftar",
			})
			return
		}

		if err.Error() == "invalid disaster status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status bencana tidak valid, hanya bisa 'active', 'resolved', atau 'archived'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat laporan bencana",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil membuat laporan bencana",
	})
}

// UpdateDisaster godoc
// @Summary Update disaster
// @Description Mengupdate laporan bencana
// @Tags Disaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Disaster ID"
// @Param disaster body structs.DisasterInput true "Disaster object"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 401 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /disasters/{id} [put]
func UpdateDisaster(c *gin.Context) {
	var input structs.DisasterInput
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

	disaster := structs.Disaster{
		ID:          id,
		Type:        input.Type,
		Location:    input.Location,
		Description: input.Description,
		Status:      input.Status,
		ReportedBy:  input.ReportedBy,
	}

	err = repository.UpdateDisaster(database.DbConnection, disaster)
	if err != nil {
		if err.Error() == "invalid disaster status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status bencana tidak valid, hanya bisa 'active', 'resolved', atau 'archived'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate laporan bencana",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengupdate laporan bencana",
	})
}

// DeleteDisaster godoc
// @Summary Delete disaster
// @Description Menghapus laporan bencana
// @Tags Disaster
// @Produce json
// @Param id path int true "Disaster ID"
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /disasters/{id} [delete]
func DeleteDisaster(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteDisaster(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Laporan bencana tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus laporan bencana",
	})
}