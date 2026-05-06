package handlers

import (
	"net/http"
	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"strconv"
	"github.com/gin-gonic/gin"
)

type SyncResponse struct {
	Chambres     []logic.Chambre     `json:"chambres"`
	Reservations []logic.Reservation `json:"reservations"`
}

func SyncHandler(c *gin.Context) {
	sinceStr := c.Query("since")
	var since int64
	if sinceStr != "" {
		val, err := strconv.ParseInt(sinceStr, 10, 64)
		if err == nil {
			since = val
		}
	}

	var chambres []logic.Chambre
	db.GetDB().Where("updated_at > ?", since).Find(&chambres)

	var reservations []logic.Reservation
	db.GetDB().Where("updated_at > ?", since).Find(&reservations)

	c.JSON(http.StatusOK, SyncResponse{
		Chambres:     chambres,
		Reservations: reservations,
	})
}
