package middleware

import (
    "net/http"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "non autorisé"})
            return
        }

        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "accès interdit"})
    }
}