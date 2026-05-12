package logic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

var (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
	RoleClient         = "client"
	RoleEmploye        = "employe"
	RoleReceptionniste = "receptionniste"
	RoleManager        = "manager"
)

func InitKeys() error {
	keyFile := os.Getenv("JWT_PRIVATE_KEY")
	if keyFile != "" {
		data, err := os.ReadFile(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read private key: %w", err)
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		publicKey = &privateKey.PublicKey
		slog.Info("RSA keys loaded from file", "file", keyFile)
		return nil
	}

	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}
	publicKey = &privateKey.PublicKey
	slog.Info("RSA key pair generated")
	return nil
}

type RefreshTokenStore struct {
	db *gorm.DB
}

func NewRefreshTokenStore(db *gorm.DB) *RefreshTokenStore {
	return &RefreshTokenStore{db: db}
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"type:uuid;index"`
	TokenHash string `gorm:"uniqueIndex"`
	ExpiresAt int64
	Revoked   bool `gorm:"default:false"`
	CreatedAt int64
}

func (s *RefreshTokenStore) AutoMigrate() error {
	return s.db.AutoMigrate(&RefreshToken{})
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *RefreshTokenStore) Create(userID string) (string, error) {
	token := generateRandomString(32)
	record := RefreshToken{
		ID:        generateRandomString(16),
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
		CreatedAt: time.Now().Unix(),
	}
	if err := s.db.Create(&record).Error; err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}
	return token, nil
}

func (s *RefreshTokenStore) Validate(token string) (string, error) {
	var record RefreshToken
	if err := s.db.Where("token_hash = ? AND revoked = ? AND expires_at > ?",
		hashToken(token), false, time.Now().Unix()).First(&record).Error; err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}
	return record.UserID, nil
}

func (s *RefreshTokenStore) Revoke(token string) error {
	return s.db.Model(&RefreshToken{}).Where("token_hash = ?", hashToken(token)).
		Update("revoked", true).Error
}

func (s *RefreshTokenStore) RevokeAllForUser(userID string) error {
	return s.db.Model(&RefreshToken{}).Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (u User) HasPermission(requiredRoles ...string) bool {
	for _, role := range requiredRoles {
		if u.Role == strings.TrimSpace(role) {
			return true
		}
	}
	return false
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func generateAccessToken(userID, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"iat":     now.Unix(),
		"exp":     now.Add(accessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func ValidateAccessToken(tokenString string) (userID, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid claims")
	}
	if claims["type"] != "access" {
		return "", "", fmt.Errorf("invalid token type")
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
	return db.Create(user).Error
}

func AuthenticateUser(db *gorm.DB, email, password string, tokenStore *RefreshTokenStore) (TokenPair, string, User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return TokenPair{}, "", User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return TokenPair{}, "", User{}, err
	}

	accessToken, err := generateAccessToken(user.ID, user.Role)
	if err != nil {
		return TokenPair{}, "", User{}, err
	}

	refreshToken, err := tokenStore.Create(user.ID)
	if err != nil {
		return TokenPair{}, "", User{}, err
	}

	user.Password = ""
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, user.Role, user, nil
}

func RefreshAccessToken(refreshToken string, tokenStore *RefreshTokenStore, db *gorm.DB) (*TokenPair, error) {
	userID, err := tokenStore.Validate(refreshToken)
	if err != nil {
		return nil, err
	}

	var user User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := tokenStore.Revoke(refreshToken); err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := tokenStore.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}
