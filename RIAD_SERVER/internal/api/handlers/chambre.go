package handlers

import (
	"net/http"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
)

func GetChambres(c *gin.Context) {
	chambres, err := logic.GetChambres(db.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chambres)
}

func CreateChambre(c *gin.Context) {
	// ... existing code ...
}

func UpdateCleaningStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		CleaningStatus string `json:"cleaning_status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chambre, err := logic.UpdateCleaningStatus(db.GetDB(), id, input.CleaningStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chambre)
}