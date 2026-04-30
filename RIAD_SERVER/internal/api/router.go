package api

import (
    "github.com/gin-gonic/gin"
    "RIAD_SERVER/internal/api/handlers"
    "RIAD_SERVER/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
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
    }

    // Protected routes
    protected := r.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/chambres", handlers.GetChambres)
        protected.POST("/chambres", middleware.RoleMiddleware("manager"), handlers.CreateChambre)
        protected.GET("/reservations", middleware.RoleMiddleware("manager", "receptionniste"), handlers.GetReservations)
        protected.POST("/reservations", handlers.CreateReservation)
        protected.PATCH("/reservations/:id/checkin", middleware.RoleMiddleware("manager", "receptionniste"), handlers.Checkin)
        protected.PATCH("/reservations/:id/checkout", middleware.RoleMiddleware("manager", "receptionniste"), handlers.Checkout)
        protected.GET("/auth/me", handlers.GetCurrentUser)
    }

    return r
}