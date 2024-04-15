package helpers

import (
	"testing"
	"time"

	_ "github.com/essaubaid/my_first_go_project/testing_config"

	"github.com/golang-jwt/jwt"
)

// GenerateTokenForTest is a helper to generate tokens with adjustable expiration
func GenerateTokenForTest(email, firstName, lastName, uid string, duration time.Duration) (string, string, error) {

	// Use the same function or directly manipulate for different scenarios
	signedToken, signedRefreshToken, err := GenerateAllTokens(email, firstName, lastName, uid)
	if err != nil {
		return "", "", err
	}
	return signedToken, signedRefreshToken, nil
}

func TestGenerateAllTokens(t *testing.T) {
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"
	uid := "123456"

	token, refreshToken, err := GenerateAllTokens(email, firstName, lastName, uid)
	if err != nil {
		t.Fatalf("GenerateAllTokens() error = %v, wantErr %v", err, false)
	}

	if token == "" || refreshToken == "" {
		t.Error("GenerateAllTokens() returned empty tokens, expected non-empty tokens")
	}

	// Parse the token to check if the claims are correct
	tokenClaims, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !tokenClaims.Valid {
		t.Error("Token is not valid")
	}

	claims, ok := tokenClaims.Claims.(*SignedDetails)
	if !ok {
		t.Fatal("Failed to parse claims as *SignedDetails")
	}

	if claims.Email != email || claims.First_name != firstName || claims.Last_name != lastName || claims.Uid != uid {
		t.Error("Claims do not match the input parameters")
	}

	expectedExpiry := time.Now().Local().Add(time.Hour * 24).Unix()
	if claims.StandardClaims.ExpiresAt != expectedExpiry {
		t.Errorf("Expected token expiration to be %v, got %v", expectedExpiry, claims.StandardClaims.ExpiresAt)
	}
}

func TestValidateToken(t *testing.T) {
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"
	uid := "123456"

	token, _, err := GenerateAllTokens(email, firstName, lastName, uid)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}

	claim, msg := ValidateToken(token)
	if claim == nil || msg != "" {
		t.Errorf("ValidateToken failed for valid token: got nil claims or non-empty msg")
	}
}
