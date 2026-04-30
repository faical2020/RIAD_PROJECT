package handlers

import (
    "net/http"
    "os"
    "RIAD_SERVER/internal/db"
    "RIAD_SERVER/internal/logic"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c *gin.Context) {
    var user logic.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)

    if err := db.GetDB().Create(&user).Error; err != nil {
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

    var user logic.User
    if err := db.GetDB().Where("email = ?", credentials.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
    })
    tokenString, _ := token.SignedString(jwtSecret)

    c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role})
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