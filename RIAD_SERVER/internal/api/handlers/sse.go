package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"RIAD_SERVER/internal/sync"
	"github.com/gin-gonic/gin"
)

func SseSyncHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		log.Println("[SSE] Connection rejected: no token in query")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required in query string"})
		return
	}
	
	log.Printf("[SSE] Web client connected to sync stream. Token: %s...", token[:10])

	eventChan := sync.GlobalBus.Subscribe()

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		if w == nil {
			return false
		}

		// Send welcome message
		fmt.Fprintf(w, "data: {\"type\": \"WELCOME\", \"id\": \"server\"}\n\n")
		w.(http.Flusher).Flush()

		// We can't easily block in a loop here without risking the context
		// But for SSE, we can check the channel.
		// However, c.Stream is called repeatedly.
		// The proper way to do this with a channel in Gin is to use a loop inside.
		
		for {
			select {
			case <-c.Request.Context().Done():
				return false
			case event := <-eventChan:
				log.Printf("[SSE] Publishing event to web client: %s for %s", event.Type, event.EntityID)
				msg := fmt.Sprintf(`{"type": "%s", "id": "%s"}`, event.Type, event.EntityID)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.(http.Flusher).Flush()
			}
		}
	})
}
