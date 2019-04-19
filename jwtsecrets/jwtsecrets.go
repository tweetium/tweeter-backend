package jwtsecrets

import (
	"github.com/sirupsen/logrus"
)

// CurrentKey is the secret key used to encode new signed tokens
var CurrentKey string

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
	CurrentKey = currentKey
}

// Clear is used for testing and mimics that initialize is not called from a previous test
func Clear() {
	keyToSecretMap = nil
	CurrentKey = ""
}

// GetSecretFor returns the secret for the given key if it exists
func GetSecretFor(keyID string) (secret []byte, ok bool) {
	secret, ok = keyToSecretMap[keyID]
	return
}

var keyToSecretMap map[string][]byte
