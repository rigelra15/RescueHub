package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateShelter godoc
// @Summary Create a shelter
// @Description Membuat shelter baru
// @Tags Shelter
// @Accept json
// @Produce json
// @Param input body structs.ShelterInput true "Data shelter"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /shelters [post]
func CreateShelter(c *gin.Context) {
	var input structs.ShelterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	shelter := structs.Shelter{
		Name:        input.Name,
		Location:    input.Location,
		CapacityTotal: input.CapacityTotal,
		CapacityRemaining: input.CapacityTotal,
		EmergencyNeeds: input.EmergencyNeeds,
		DisasterID: input.DisasterID,
	}

	err := repository.CreateShelter(database.DbConnection, shelter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat shelter",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Shelter berhasil dibuat",
	})
}

// GetAllShelters godoc
// @Summary Get all shelters
// @Description Mendapatkan daftar shelter
// @Tags Shelter
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /shelters [get]
func GetAllShelters(c *gin.Context) {
	shelters, err := repository.GetAllShelters(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar shelters",
		})
		return
	}

	if len(shelters) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tidak ada shelter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": shelters,
	})
}

// GetShelterByID godoc
// @Summary Get a shelter by ID
// @Description Mendapatkan shelter berdasarkan ID
// @Tags Shelter
// @Accept json
// @Produce json
// @Param id path int true "ID Shelter"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /shelters/{id} [get]
func GetShelterByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	shelter, err := repository.GetShelterByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Shelter tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": shelter,
	})
}

// UpdateShelter godoc
// @Summary Update a shelter
// @Description Memperbarui shelter
// @Tags Shelter
// @Accept json
// @Produce json
// @Param id path int true "ID Shelter"
// @Param input body structs.ShelterInput true "Data shelter"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /shelters/{id} [put]
func UpdateShelter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.Shelter
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}
	input.ID = id

	err = repository.UpdateShelter(database.DbConnection, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate shelter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shelter berhasil diperbarui",
	})
}

// DeleteShelter godoc
// @Summary Delete a shelter
// @Description Menghapus shelter
// @Tags Shelter
// @Accept json
// @Produce json
// @Param id path int true "ID Shelter"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /shelters/{id} [delete]
func DeleteShelter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteShelter(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus shelter",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shelter berhasil dihapus",
	})
}
