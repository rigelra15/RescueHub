package controllers

import (
	"RescueHub/database"
	"RescueHub/repository"
	"RescueHub/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEvacuationRoute godoc
// @Summary Create an evacuation route
// @Description Mencatat jalur evakuasi
// @Tags EvacuationRoute
// @Accept json
// @Produce json
// @Param input body structs.EvacuationRouteInput true "Input data jalur evakuasi"
// @Success 201 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /evacuation_routes [post]
func CreateEvacuationRoute(c *gin.Context) {
	var input structs.EvacuationRouteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	evacuationRoute := &structs.EvacuationRoute{
		DisasterID: 	input.DisasterID,
		Origin:     	input.Origin,
		Destination: 	input.Destination,
		Distance:  		input.Distance,
		Route:      	input.Route,
		Status:     	input.Status,
	}	

	err := repository.CreateEvacuationRoute(database.DbConnection, evacuationRoute)
	if err != nil {
		if err.Error() == "invalid evacuation route status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status jalur evakuasi tidak valid, hanya bisa 'safe', 'risky', atau 'blocked'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mencatat jalur evakuasi",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jalur evakuasi berhasil dicatat",
		"result":  evacuationRoute,
	})
}

// GetAllEvacuationRoutes godoc
// @Summary Get all evacuation routes
// @Description Mendapatkan daftar jalur evakuasi
// @Tags EvacuationRoute
// @Accept json
// @Produce json
// @Success 200 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Router /evacuation_routes [get]
func GetAllEvacuationRoutes(c *gin.Context) {
	routes, err := repository.GetAllEvacuationRoutes(database.DbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mendapatkan daftar jalur evakuasi",
		})
		return
	}

	if len(routes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada daftar jalur evakuasi yang tersedia",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": routes,
	})
}

// GetEvacuationRouteByID godoc
// @Summary Get an evacuation route by ID
// @Description Mendapatkan jalur evakuasi berdasarkan ID
// @Tags EvacuationRoute
// @Accept json
// @Produce json
// @Param id path int true "ID jalur evakuasi"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 404 {object} structs.APIResponse
// @Router /evacuation_routes/{id} [get]
func GetEvacuationRouteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	route, err := repository.GetEvacuationRouteByID(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Jalur evakuasi tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": route,
	})
}

// UpdateEvacuationRoute godoc
// @Summary Update an evacuation route
// @Description Memperbarui jalur evakuasi
// @Tags EvacuationRoute
// @Accept json
// @Produce json
// @Param id path int true "ID jalur evakuasi"
// @Param input body structs.EvacuationRouteInput true "Input data jalur evakuasi"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /evacuation_routes/{id} [put]
func UpdateEvacuationRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var input structs.EvacuationRouteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	evacuationRoute := structs.EvacuationRoute{
		ID:          id,
		DisasterID:  input.DisasterID,
		Origin:      input.Origin,
		Destination: input.Destination,
		Distance:    input.Distance,
		Route:       input.Route,
		Status:      input.Status,
	}

	err = repository.UpdateEvacuationRoute(database.DbConnection, evacuationRoute)
	if err != nil {
		if err.Error() == "invalid evacuation route status" {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "Status jalur evakuasi tidak valid, hanya bisa 'safe', 'risky', atau 'blocked'",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate jalur evakuasi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jalur evakuasi berhasil diperbarui",
	})
}

// DeleteEvacuationRoute godoc
// @Summary Delete an evacuation route
// @Description Menghapus jalur evakuasi
// @Tags EvacuationRoute
// @Accept json
// @Produce json
// @Param id path int true "ID jalur evakuasi"
// @Success 200 {object} structs.APIResponse
// @Failure 400 {object} structs.APIResponse
// @Failure 500 {object} structs.APIResponse
// @Security BearerAuth
// @Router /evacuation_routes/{id} [delete]
func DeleteEvacuationRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	err = repository.DeleteEvacuationRoute(database.DbConnection, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus jalur evakuasi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jalur evakuasi berhasil dihapus",
	})
}
