package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"tweeter/db/models/user"
	"tweeter/jwtsecrets"
)

// CookieName is the name of the cookie stored
var CookieName = "user_token"

// ExpirationDuration is the default expiration duration for the claim
var ExpirationDuration = time.Minute * 30

// Claims is the exposed claims stored in the token
type Claims struct {
	UserID user.ID
}

// GenerateToken generates a signed token for the claims given
func GenerateToken(claims Claims) (string, error) {
	expirationTime := time.Now().Add(ExpirationDuration)
	return GenerateTokenWithExpiration(claims, expirationTime)
}

// GenerateTokenWithExpiration generates a signed token with the expiration time given
func GenerateTokenWithExpiration(claims Claims, expirationTime time.Time) (string, error) {
	return jwtsecrets.NewSignedClaims(&parsedClaims{
		UserID: claims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})
}

// ParseToken validates a token string and returns the claims, otherwise returns the error
func ParseToken(rawToken string) (Claims, error) {
	var parsed parsedClaims
	err := jwtsecrets.ParseAndVerifyClaims(rawToken, &parsed)

	if err != nil {
		return Claims{}, err
	}

	claims := Claims{
		UserID: parsed.UserID,
	}
	return claims, err
}

type parsedClaims struct {
	UserID user.ID `json:"user_id"`
	jwt.StandardClaims
}
