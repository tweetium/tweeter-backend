package jwt

import (
	"errors"
	"fmt"
	"time"
	"tweeter/db/models/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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

// InitializeWithSecretsMap parses a map of key id -> jwt secret key
// The currentID is used for all new tokens signed
func InitializeWithSecretsMap(idToSecretStringMap map[string]string, currentID string) {
	idToSecretMap = map[string][]byte{}
	for id, secret := range idToSecretStringMap {
		idToSecretMap[id] = []byte(secret)
	}

	if _, ok := idToSecretMap[currentID]; !ok {
		logrus.WithFields(logrus.Fields{
			"secretsMap": idToSecretStringMap,
			"currentID":  currentID,
		}).Panic("Current secret ID not in secrets map")
	}
	currentKeyID = currentID
}

// ClearSecretsMap is used for testing and mimics that initialize is not called from a previous test
func ClearSecretsMap() {
	idToSecretMap = nil
	currentKeyID = ""
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
	token.Header["kid"] = currentKeyID
	secret, ok := getSecretFor(currentKeyID)
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
		secret, ok = getSecretFor(keyIDString)
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

var idToSecretMap map[string][]byte
var currentKeyID string

type parsedClaims struct {
	UserID user.ID `json:"user_id"`
	jwt.StandardClaims
}

func getSecretFor(keyID string) (secret []byte, ok bool) {
	secret, ok = idToSecretMap[keyID]
	return
}
