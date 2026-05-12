package handlers

import (
	"net/http"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tokenStore = logic.NewRefreshTokenStore(nil)

func SetTokenStore(store *logic.RefreshTokenStore) {
	tokenStore = store
}

func Register(c *gin.Context) {
	var user logic.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := logic.RegisterUser(db.GetDB(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur création"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "utilisateur créé"})
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, role, user, err := logic.AuthenticateUser(db.GetDB(), credentials.Email, credentials.Password, tokenStore)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"role":          role,
		"user":          user,
	})
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token requis"})
		return
	}

	tokenPair, err := logic.RefreshAccessToken(req.RefreshToken, tokenStore, db.GetDB())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token de rafraîchissement invalide ou expiré"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

func Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req)

	if req.RefreshToken != "" {
		tokenStore.Revoke(req.RefreshToken)
	}

	userID, _ := c.Get("user_id")
	if uid, ok := userID.(string); ok {
		tokenStore.RevokeAllForUser(uid)
	}

	c.JSON(http.StatusOK, gin.H{"message": "déconnexion réussie"})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user logic.User
	if err := db.GetDB().First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "utilisateur non trouvé"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur serveur"})
		return
	}
	c.JSON(http.StatusOK, user)
}
