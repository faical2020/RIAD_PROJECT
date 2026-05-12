package logic

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const (
	RoleClient         = "client"
	RoleEmploye        = "employe"
	RoleReceptionniste = "receptionniste"
	RoleManager        = "manager"
)

func (u User) HasPermission(requiredRoles ...string) bool {
	for _, role := range requiredRoles {
		if u.Role == strings.TrimSpace(role) {
			return true
		}
	}
	return false
}

func ValidateToken(tokenString string) (userID, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil
	}
	userID, _ = claims["user_id"].(string)
	role, _ = claims["role"].(string)
	return userID, role, nil
}

func RegisterUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func AuthenticateUser(db *gorm.DB, email, password string) (string, string, User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", User{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", User{}, err
	}

	user.Password = ""
	return tokenString, user.Role, user, nil
}