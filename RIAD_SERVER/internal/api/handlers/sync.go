package handlers

import (
	"net/http"
	"strconv"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
)

func SyncHandler(c *gin.Context) {
	sinceStr := c.Query("since")
	var since int64
	if sinceStr != "" {
		val, err := strconv.ParseInt(sinceStr, 10, 64)
		if err == nil {
			since = val
		}
	}

	data, err := logic.GetSyncUpdates(db.GetDB(), since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
