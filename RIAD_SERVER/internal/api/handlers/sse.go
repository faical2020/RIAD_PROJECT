package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"RIAD_SERVER/internal/eventbus"
	"github.com/gin-gonic/gin"
)

func SseSyncHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		slog.Warn("SSE connection rejected: no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required in query string"})
		return
	}

	slog.Info("SSE web client connected", "token_prefix", token[:10])

	eventChan := eventbus.GlobalBus.Subscribe()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		if w == nil {
			return false
		}

		fmt.Fprintf(w, "data: {\"type\": \"WELCOME\", \"id\": \"server\"}\n\n")
		w.(http.Flusher).Flush()

		for {
			select {
			case <-c.Request.Context().Done():
				return false
			case event := <-eventChan:
				slog.Debug("SSE publishing event", "type", event.Type, "id", event.EntityID)
				msg := fmt.Sprintf(`{"type": "%s", "id": "%s"}`, event.Type, event.EntityID)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.(http.Flusher).Flush()
			}
		}
	})
}
