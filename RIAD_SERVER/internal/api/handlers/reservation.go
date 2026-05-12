package handlers

import (
	"net/http"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
)

func GetReservations(c *gin.Context) {
	reservations, err := logic.GetReservations(db.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

func CreateReservation(c *gin.Context) {
	var res logic.Reservation
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := logic.CreateReservation(db.GetDB(), &res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func Checkin(c *gin.Context) {
	id := c.Param("id")

	res, err := logic.CheckinReservation(db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func Checkout(c *gin.Context) {
	id := c.Param("id")

	res, err := logic.CheckoutReservation(db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
