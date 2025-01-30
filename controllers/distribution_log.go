package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDistributionLog godoc
// @Summary Create a distribution log
// @Description Mencatat distribusi bantuan
// @Tags DistributionLog
// @Accept json
// @Produce json
// @Param input body structs.DistributionLogInput true "Data distribusi bantuan"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /distribution_logs [post]
func CreateDistributionLog(c *gin.Context) {
	var input structs.DistributionLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	distributionLog := structs.DistributionLog{
		LogisticID: input.LogisticID,
		Origin: 	 input.Origin,
		Destination: input.Destination,
		Distance: input.Distance,
		SenderName: input.SenderName,
		RecipientName: input.RecipientName,
		QuantitySent: input.QuantitySent,
		SentAt: input.SentAt,
	}

	err := repository.CreateDistributionLog(database.DbConnection, distributionLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat distribusi bantuan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Distribusi bantuan berhasil dicatat",
	})
}

// GetAllDistributionLogs godoc
// @Summary Get all distribution logs
// @Description Mendapatkan daftar distribusi bantuan
// @Tags DistributionLog
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /distribution_logs [get]
func GetAllDistributionLogs(c *gin.Context) {
	logs, err := repository.GetAllDistributionLogs(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar distribusi bantuan",
		})
		return
	}

	if len(logs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar distribusi bantuan yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": logs,
	})
}

// GetDistributionLogByID godoc
// @Summary Get distribution log by ID
// @Description Mendapatkan distribusi bantuan berdasarkan ID
// @Tags DistributionLog
// @Accept json
// @Produce json
// @Param id path int true "Distribution Log ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /distribution_logs/{id} [get]
func GetDistributionLogByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	log, err := repository.GetDistributionLogByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Distribusi bantuan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": log,
	})
}

// UpdateDistributionLog godoc
// @Summary Update a distribution log
// @Description Memperbarui distribusi bantuan
// @Tags DistributionLog
// @Accept json
// @Produce json
// @Param id path int true "Distribution Log ID"
// @Param input body structs.DistributionLogInput true "Data distribusi bantuan"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /distribution_logs/{id} [put]
func UpdateDistributionLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.DistributionLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	distributionLog := structs.DistributionLog{
		ID:         id,
		LogisticID: input.LogisticID,
		Origin:     input.Origin,
		Destination: input.Destination,
		Distance: input.Distance,
		SenderName: input.SenderName,
		RecipientName: input.RecipientName,
		QuantitySent: input.QuantitySent,
		SentAt: input.SentAt,
	}

	err = repository.UpdateDistributionLog(database.DbConnection, distributionLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate distribusi bantuan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Distribusi bantuan berhasil diperbarui",
	})
}

// DeleteDistributionLog godoc
// @Summary Delete a distribution log
// @Description Menghapus distribusi bantuan
// @Tags DistributionLog
// @Accept json
// @Produce json
// @Param id path int true "Distribution Log ID"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /distribution_logs/{id} [delete]
func DeleteDistributionLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteDistributionLog(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus distribusi bantuan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Distribusi bantuan berhasil dihapus",
	})
}
