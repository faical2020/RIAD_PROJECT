package handlers

import (
	"net/http"
	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"RIAD_SERVER/internal/sync"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"github.com/gin-gonic/gin"
)


func GetChambres(c *gin.Context) {
    var chambres []logic.Chambre
    db.GetDB().Find(&chambres)
    c.JSON(http.StatusOK, chambres)
}

func CreateChambre(c *gin.Context) {
    var chambre logic.Chambre
    if err := c.ShouldBindJSON(&chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := logic.ValidateChambre(chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	db.GetDB().Create(&chambre)
	
	sync.GlobalBus.Publish("ROOM_UPDATED", chambre.ID, &pb.Room{
		Id:          chambre.ID,
		Number:      int32(chambre.Numero),
		Type:        chambre.Type,
		Price:       chambre.Prix,
		Description: chambre.Description,
		Equipments:  chambre.Equipements,
		Status:      chambre.Statut,
	})

	c.JSON(http.StatusCreated, chambre)

}