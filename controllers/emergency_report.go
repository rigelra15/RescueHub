package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEmergencyReport godoc
// @Summary Create an emergency report
// @Description Mencatat laporan darurat
// @Tags EmergencyReport
// @Accept json
// @Produce json
// @Param input body structs.EmergencyReportInput true "Input data laporan darurat"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /emergency_reports [post]
func CreateEmergencyReport(c *gin.Context) {
	var input structs.EmergencyReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	emergencyReport := structs.EmergencyReport{
		UserID:      input.UserID,
		DisasterID:  input.DisasterID,
		Description: input.Description,
		Location:    input.Location,
	}

	err := repository.CreateEmergencyReport(database.DbConnection, emergencyReport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat laporan darurat",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Laporan darurat berhasil dicatat",
	})
}

// GetAllEmergencyReports godoc
// @Summary Get all emergency reports
// @Description Mendapatkan daftar laporan darurat
// @Tags EmergencyReport
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /emergency_reports [get]
func GetAllEmergencyReports(c *gin.Context) {
	reports, err := repository.GetAllEmergencyReports(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar laporan darurat",
		})
		return
	}

	if len(reports) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar laporan darurat yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": reports,
	})
}

// GetEmergencyReportByID godoc
// @Summary Get an emergency report by ID
// @Description Mendapatkan laporan darurat berdasarkan ID
// @Tags EmergencyReport
// @Accept json
// @Produce json
// @Param id path int true "ID laporan darurat"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /emergency_reports/{id} [get]
func GetEmergencyReportByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	report, err := repository.GetEmergencyReportByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Laporan darurat tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": report,
	})
}

// UpdateEmergencyReport godoc
// @Summary Update an emergency report
// @Description Memperbarui laporan darurat
// @Tags EmergencyReport
// @Accept json
// @Produce json
// @Param id path int true "ID laporan darurat"
// @Param input body structs.EmergencyReportInput true "Input data laporan darurat"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /emergency_reports/{id} [put]
func UpdateEmergencyReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.EmergencyReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}
	
	emergencyReport := structs.EmergencyReport{
		ID:          id,
		UserID:      input.UserID,
		DisasterID:  input.DisasterID,
		Description: input.Description,
		Location:    input.Location,
	}

	err = repository.UpdateEmergencyReport(database.DbConnection, emergencyReport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate laporan darurat",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan darurat berhasil diperbarui",
	})
}

// DeleteEmergencyReport godoc
// @Summary Delete an emergency report
// @Description Menghapus laporan darurat
// @Tags EmergencyReport
// @Accept json
// @Produce json
// @Param id path int true "ID laporan darurat"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /emergency_reports/{id} [delete]
func DeleteEmergencyReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteEmergencyReport(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus laporan darurat",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan darurat berhasil dihapus",
	})
}
