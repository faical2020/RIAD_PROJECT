package handlers

import (
    "net/http"
    "RIAD_SERVER/internal/db"
    "RIAD_SERVER/internal/logic"
    "github.com/gin-gonic/gin"
)

func GetReservations(c *gin.Context) {
    var reservations []logic.Reservation
    db.GetDB().Find(&reservations)
    c.JSON(http.StatusOK, reservations)
}

func CreateReservation(c *gin.Context) {
    var res logic.Reservation
    if err := c.ShouldBindJSON(&res); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var chambre logic.Chambre
    if err := db.GetDB().First(&chambre, "id = ?", res.ChambreID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "chambre non trouvée"})
        return
    }

    if err := logic.ValidateReservation(res, chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.GetDB().Create(&res)
    c.JSON(http.StatusCreated, res)
}

func Checkin(c *gin.Context) {
    id := c.Param("id")
    var res logic.Reservation
    if err := db.GetDB().First(&res, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "réservation non trouvée"})
        return
    }

    var chambre logic.Chambre
    db.GetDB().First(&chambre, "id = ?", res.ChambreID)

    if err := res.Checkin(&chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.GetDB().Save(&res).Save(&chambre)
    c.JSON(http.StatusOK, res)
}

func Checkout(c *gin.Context) {
    id := c.Param("id")
    var res logic.Reservation
    if err := db.GetDB().First(&res, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "réservation non trouvée"})
        return
    }

    var chambre logic.Chambre
    db.GetDB().First(&chambre, "id = ?", res.ChambreID)

    if err := res.Checkout(&chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.GetDB().Save(&res).Save(&chambre)
    c.JSON(http.StatusOK, res)
}