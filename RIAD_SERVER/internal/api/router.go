package api

import (
    "github.com/gin-gonic/gin"
    "RIAD_SERVER/internal/api/handlers"
    "RIAD_SERVER/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware - Application stricte et globale
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/refresh", handlers.RefreshToken)
		public.GET("/sync/events", handlers.SseSyncHandler)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/auth/logout", handlers.Logout)
		protected.GET("/sync", handlers.SyncHandler)
		protected.GET("/chambres", handlers.GetChambres)
		protected.POST("/chambres", middleware.RoleMiddleware("manager"), handlers.CreateChambre)
		protected.PATCH("/chambres/:id/cleaning", handlers.UpdateCleaningStatus)
		protected.GET("/reservations/mine", handlers.GetMyReservations)
		protected.GET("/reservations", middleware.RoleMiddleware("manager", "receptionniste"), handlers.GetReservations)
		protected.POST("/reservations", handlers.CreateReservation)
		protected.PATCH("/reservations/:id/checkin", middleware.RoleMiddleware("manager", "receptionniste"), handlers.Checkin)
		protected.PATCH("/reservations/:id/checkout", middleware.RoleMiddleware("manager", "receptionniste"), handlers.Checkout)
		protected.GET("/auth/me", handlers.GetCurrentUser)

		// Services & Consommations
		protected.GET("/services", handlers.GetServices)
		protected.POST("/services", middleware.RoleMiddleware("manager"), handlers.CreateService)
		protected.PUT("/services/:id", middleware.RoleMiddleware("manager"), handlers.UpdateService)
		protected.DELETE("/services/:id", middleware.RoleMiddleware("manager"), handlers.DeleteService)
		protected.GET("/reservations/:id/consommations", handlers.GetConsommations)
		protected.POST("/reservations/:id/consommations", handlers.AddConsommation)
		protected.DELETE("/consommations/:id", middleware.RoleMiddleware("manager"), handlers.DeleteConsommation)
		protected.GET("/reservations/:id/facture", handlers.GetFacture)
		protected.POST("/reservations/:id/paiement", handlers.AddPaiement)
	}

	return r
}
