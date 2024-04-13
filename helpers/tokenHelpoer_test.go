package helpers

import (
	"testing"
	"time"

	_ "github.com/essaubaid/my_first_go_project/testing_config"

	"github.com/golang-jwt/jwt"
)

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
