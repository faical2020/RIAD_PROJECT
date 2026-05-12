package handlers

import (
	"net/http"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	token, role, user, err := logic.AuthenticateUser(db.GetDB(), credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "role": role, "user": user})
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