package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"tweeter/db/models/user"
	"tweeter/jwtsecrets"
)

// CookieName is the name of the cookie stored
var CookieName = "user_token"

// ExpirationDuration is the default expiration duration for the claim
var ExpirationDuration = time.Minute * 30

// ErrTokenExpired is the error returned when the token is expired
var ErrTokenExpired = errors.New("Token is expired")

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &parsedClaims{
		UserID: claims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})
	// set the key used to sign this token (kid == keyID)
	token.Header["kid"] = jwtsecrets.CurrentKey
	secret, ok := jwtsecrets.GetSecretFor(jwtsecrets.CurrentKey)
	if !ok {
		return "", errors.New("Current secret key invalid - was usersJWT initialized?")
	}

	return token.SignedString(secret)
}

// ParseToken validates a token string and returns the claims, otherwise returns the error
func ParseToken(rawToken string) (Claims, error) {
	var parsed parsedClaims
	token, err := jwt.ParseWithClaims(rawToken, &parsed, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		keyID, ok := token.Header["kid"]
		if !ok {
			return nil, errors.New("Missing 'kid' (keyID) field for usersJWT token")
		}

		var keyIDString string
		keyIDString, ok = keyID.(string)
		if !ok {
			return nil, errors.New("keyID for usersJWT token was not string")
		}

		var secret []byte
		secret, ok = jwtsecrets.GetSecretFor(keyIDString)
		if !ok {
			return nil, fmt.Errorf("Failed to find secret for keyID: %v", keyID)
		}
		return secret, nil
	})

	// Return the error first before checking token.Valid
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors == jwt.ValidationErrorExpired {
				return Claims{}, ErrTokenExpired
			}
		}

		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("Token is not valid")
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
