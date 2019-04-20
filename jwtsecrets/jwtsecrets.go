package jwtsecrets

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// ErrTokenExpired is the error returned when the token is expired
var ErrTokenExpired = errors.New("Token is expired")

// Init parses a map of key id -> jwt secret key
// The currentKey is used for all new tokens signed
func Init(keyToSecretStringMap map[string]string, currentKey string) {
	keyToSecretMap = map[string][]byte{}
	for id, secret := range keyToSecretStringMap {
		keyToSecretMap[id] = []byte(secret)
	}

	if _, ok := keyToSecretMap[currentKey]; !ok {
		logrus.WithFields(logrus.Fields{
			"secretsMap": keyToSecretStringMap,
			"currentKey": currentKey,
		}).Panic("Current secret ID not in secrets map")
	}
	currentSecretKey = currentKey
}

// Clear is used for testing and mimics that initialize is not called from a previous test
func Clear() {
	keyToSecretMap = nil
	currentSecretKey = ""
}

// NewSignedClaims returns a signed string representing the signed claims
// Any error returned means jwtsecrets was initialized improperly and cannot be fixed
func NewSignedClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set the key used to sign this token (kid == keyID)
	token.Header["kid"] = currentSecretKey
	secret, ok := getSecretFor(currentSecretKey)
	if !ok {
		return "", errors.New("Current secret key invalid - was usersJWT initialized?")
	}

	return token.SignedString(secret)
}

// ParseAndVerifyClaims parses a token string and verifies the secret
func ParseAndVerifyClaims(rawToken string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
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
				return ErrTokenExpired
			}
		}
	}

	return err
}

// getSecretFor returns the secret for the given key if it exists
func getSecretFor(keyID string) (secret []byte, ok bool) {
	secret, ok = keyToSecretMap[keyID]
	return
}

var keyToSecretMap map[string][]byte
var currentSecretKey string
