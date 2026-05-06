package handlers

import (
	"net/http"
	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"RIAD_SERVER/internal/sync"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"github.com/gin-gonic/gin"
	"time"
)


func GetChambres(c *gin.Context) {
    var chambres []logic.Chambre
    db.GetDB().Find(&chambres)
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

	var chambre logic.Chambre
	if err := db.GetDB().First(&chambre, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chambre non trouvée"})
		return
	}

	chambre.CleaningStatus = input.CleaningStatus
	chambre.UpdatedAt = time.Now().Unix()

	if err := db.GetDB().Save(&chambre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors de la mise à jour"})
		return
	}

	sync.GlobalBus.Publish("ROOM_UPDATED", chambre.ID, &pb.Room{
		Id:          chambre.ID,
		Number:      int32(chambre.Numero),
		Type:        chambre.Type,
		Price:       chambre.Prix,
		Description: chambre.Description,
		Equipments:  chambre.Equipements,
		Status:      chambre.Statut,
	})

	c.JSON(http.StatusOK, chambre)
}