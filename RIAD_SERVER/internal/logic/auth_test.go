package logic

import (
	"crypto/rand"
	"crypto/rsa"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelWarn)
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		os.Exit(1)
	}
	publicKey = &privateKey.PublicKey
	os.Exit(m.Run())
}

func TestGenerateAccessToken(t *testing.T) {
	token, err := generateAccessToken("user-123", "manager")
	if err != nil {
		t.Fatalf("generateAccessToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateAccessToken_Valid(t *testing.T) {
	token, err := generateAccessToken("user-123", "manager")
	if err != nil {
		t.Fatal(err)
	}

	userID, role, err := ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("expected user-123, got %s", userID)
	}
	if role != "manager" {
		t.Errorf("expected manager, got %s", role)
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	// Temporarily set a very short expiry
	oldTTL := accessTokenTTL
	accessTokenTTL = -1 * time.Minute
	defer func() { accessTokenTTL = oldTTL }()

	token, err := generateAccessToken("user-123", "manager")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = ValidateAccessToken(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateAccessToken_WrongKey(t *testing.T) {
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": "user-123",
		"role":    "manager",
		"type":    "access",
	})
	signed, err := token.SignedString(otherKey)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = ValidateAccessToken(signed)
	if err == nil {
		t.Error("expected error when token signed by different key")
	}
}

func TestValidateAccessToken_InvalidType(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": "user-123",
		"role":    "manager",
		"type":    "refresh",
	})
	signed, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = ValidateAccessToken(signed)
	if err == nil {
		t.Error("expected error for non-access token type")
	}
}

func TestHashToken(t *testing.T) {
	h1 := hashToken("test-token")
	h2 := hashToken("test-token")
	h3 := hashToken("other-token")

	if h1 != h2 {
		t.Error("same token should produce same hash")
	}
	if h1 == h3 {
		t.Error("different tokens should produce different hashes")
	}
	if len(h1) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(h1))
	}
}

func TestGenerateRandomString(t *testing.T) {
	s1 := generateRandomString(32)
	s2 := generateRandomString(32)

	if len(s1) != 64 {
		t.Errorf("expected 64 hex chars (32 bytes), got %d", len(s1))
	}
	if s1 == s2 {
		t.Error("random strings should be different")
	}
}

func TestUserHasPermission(t *testing.T) {
	user := User{Role: "manager"}

	if !user.HasPermission("manager") {
		t.Error("manager should have manager permission")
	}
	if !user.HasPermission("manager", "receptionniste") {
		t.Error("manager should be in [manager, receptionniste]")
	}
	if user.HasPermission("client") {
		t.Error("manager should not have client permission")
	}
}
