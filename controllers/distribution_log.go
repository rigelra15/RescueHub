package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
// @Security BearerAuth
// @Router /distribution_logs [post]
func CreateDistributionLog(c *gin.Context) {
	var input structs.DistributionLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	layout := "02/01/2006 15:04"
	parsedSentAt, err := time.Parse(layout, input.SentAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format tanggal harus 'DD/MM/YYYY HH:mm'",
		})
		return
	}

	distributionLog := &structs.DistributionLog{
		LogisticID:    input.LogisticID,
		Origin:        input.Origin,
		Destination:   input.Destination,
		Distance:      input.Distance,
		SenderName:    input.SenderName,
		RecipientName: input.RecipientName,
		QuantitySent:  input.QuantitySent,
		SentAt:        parsedSentAt,
	}

	err = repository.CreateDistributionLog(database.DbConnection, distributionLog)
	if err != nil {
		fmt.Println("Error Query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat distribusi bantuan",
		})
		return
	}

	responseSentAt := parsedSentAt.Format(layout)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Distribusi bantuan berhasil dicatat",
		"result": gin.H{
			"logistic_id":    distributionLog.LogisticID,
			"origin":         distributionLog.Origin,
			"destination":    distributionLog.Destination,
			"distance":       distributionLog.Distance,
			"sender_name":    distributionLog.SenderName,
			"recipient_name": distributionLog.RecipientName,
			"quantity_sent":  distributionLog.QuantitySent,
			"sent_at":        responseSentAt,
		},
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
// @Security BearerAuth
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
// @Security BearerAuth
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
// @Security BearerAuth
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

	layout := "02/01/2006 15:04"
	parsedSentAt, err := time.Parse(layout, input.SentAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format tanggal harus 'DD/MM/YYYY HH:mm'",
		})
		return
	}

	distributionLog := &structs.DistributionLog{
		ID:            id,
		LogisticID:    input.LogisticID,
		Origin:        input.Origin,
		Destination:   input.Destination,
		Distance:      input.Distance,
		SenderName:    input.SenderName,
		RecipientName: input.RecipientName,
		QuantitySent:  input.QuantitySent,
		SentAt:        parsedSentAt, 
	}

	err = repository.UpdateDistributionLog(database.DbConnection, *distributionLog)
	if err != nil {
		if err.Error() == "distribution log not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Log distribusi tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate distribusi bantuan",
		})
		return
	}

	responseSentAt := parsedSentAt.Format(layout)

	c.JSON(http.StatusOK, gin.H{
		"message": "Distribusi bantuan berhasil diperbarui",
		"result": gin.H{
			"id":             distributionLog.ID,
			"logistic_id":    distributionLog.LogisticID,
			"origin":         distributionLog.Origin,
			"destination":    distributionLog.Destination,
			"distance":       distributionLog.Distance,
			"sender_name":    distributionLog.SenderName,
			"recipient_name": distributionLog.RecipientName,
			"quantity_sent":  distributionLog.QuantitySent,
			"sent_at":        responseSentAt,
		},
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
// @Security BearerAuth
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
