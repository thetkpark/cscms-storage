package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const JWTSECRET = "someRandomSecretThatReallyReallyLong"

func TestJWTGeneration(t *testing.T) {
	userID := "1"
	jwtManager := NewJWTManager(JWTSECRET)
	token, err := jwtManager.Generate(userID)
	require.NoError(t, err)
	require.Greater(t, len(token), 0)
}

func TestJWTValidation(t *testing.T) {
	userID := "1"
	signedTokenString, err := generateTestJWT(userID, JWTSECRET, time.Now().Add(time.Hour))
	require.NoError(t, err)

	jwtManager := NewJWTManager(JWTSECRET)
	jwtUserID, err := jwtManager.Validate(signedTokenString)
	require.NoError(t, err)
	require.Equal(t, userID, jwtUserID, "expected", userID, "got", jwtUserID)
}

func TestExpiredJWTValidation(t *testing.T) {
	userID := "1"
	signedTokenString, err := generateTestJWT(userID, JWTSECRET, time.Now().Add(time.Hour*-1))
	require.NoError(t, err)

	jwtManager := NewJWTManager(JWTSECRET)
	jwtUserID, err := jwtManager.Validate(signedTokenString)
	require.Error(t, err)
	require.Empty(t, jwtUserID)
}

func TestInvalidJWTValidation(t *testing.T) {
	userID := "1"
	signedTokenString, err := generateTestJWT(userID, "invalidSecretKey", time.Now().Add(time.Hour))
	require.NoError(t, err)

	jwtManager := NewJWTManager(JWTSECRET)
	jwtUserID, err := jwtManager.Validate(signedTokenString)
	require.Error(t, err)
	require.Empty(t, jwtUserID)
}

func generateTestJWT(userID string, secret string, expiredAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: expiredAt.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userID,
	})
	return token.SignedString([]byte(secret))
}
