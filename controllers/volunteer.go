package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateVolunteer godoc
// @Summary Create a volunteer
// @Description Mencatat relawan
// @Tags Volunteer
// @Accept json
// @Produce json
// @Param input body structs.VolunteerInput true "Data relawan"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /volunteers [post]
func CreateVolunteer(c *gin.Context) {
	var input structs.VolunteerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	volunteer := &structs.Volunteer{
		UserID		: input.UserID,
		DisasterID: input.DisasterID,
		Skill			: input.Skill,
		Location	: input.Location,
		Status		: input.Status,
	}

	err := repository.CreateVolunteer(database.DbConnection, volunteer)
	if err != nil {
		if err.Error() == "invalid volunteer status" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status relawan tidak valid, hanya bisa 'available', 'on_mission', atau 'completed'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat relawan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Relawan berhasil dicatat",
		"result":  volunteer,
	})
}

// GetAllVolunteers godoc
// @Summary Get all volunteers
// @Description Mendapatkan daftar relawan
// @Tags Volunteer
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Security BearerAuth
// @Router /volunteers [get]
func GetAllVolunteers(c *gin.Context) {
	volunteers, err := repository.GetAllVolunteers(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar relawan yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": volunteers,
	})
}

// GetVolunteerByID godoc
// @Summary Get a volunteer by ID
// @Description Mendapatkan relawan berdasarkan ID
// @Tags Volunteer
// @Accept json
// @Produce json
// @Param id path int true "ID relawan"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Security BearerAuth
// @Router /volunteers/{id} [get]
func GetVolunteerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	volunteer, err := repository.GetVolunteerByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Relawan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": volunteer,
	})
}

// UpdateVolunteer godoc
// @Summary Update a volunteer
// @Description Memperbarui relawan
// @Tags Volunteer
// @Accept json
// @Produce json
// @Param id path int true "ID relawan"
// @Param input body structs.VolunteerInput true "Data relawan"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /volunteers/{id} [put]
func UpdateVolunteer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.VolunteerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}
	
	volunteerInput := structs.Volunteer{
		ID:       id,
		UserID:  input.UserID,
		DisasterID: input.DisasterID,
		Skill:    input.Skill,
		Location: input.Location,
		Status:   input.Status,
	}

	err = repository.UpdateVolunteer(database.DbConnection, volunteerInput)
	if err != nil {
		if err.Error() == "invalid volunteer status" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status volunteer tidak valid, hanya bisa 'available', 'on_mission', atau 'completed'",
			})
			return
		}

		if err.Error() == "volunteer not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Volunteer tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate relawan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relawan berhasil diperbarui",
	})
}

// DeleteVolunteer godoc
// @Summary Delete a volunteer
// @Description Menghapus relawan
// @Tags Volunteer
// @Accept json
// @Produce json
// @Param id path int true "ID relawan"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /volunteers/{id} [delete]
func DeleteVolunteer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteVolunteer(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus relawan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relawan berhasil dihapus",
	})
}
