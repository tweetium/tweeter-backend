package main

import (
	"encoding/json"

	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver
	"github.com/sirupsen/logrus"

	"tweeter/db"
	"tweeter/handlers"
	"tweeter/jwtsecrets"
	"tweeter/log"
	"tweeter/util"
)

func main() {
	log.Init()
	logrus.Info("Initialized logger")

	port := util.GetEnvOrDefaultUInt32("PORT", "80")
	dbURL := util.MustGetEnv("DATABASE_URL")

	initUserJWT()
	logrus.Info("Initialized User JWT store")

	err := db.Init(dbURL)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}

	handlers.RunWebserver(port)
}

func initUserJWT() {
	rawSecretsMap := util.MustGetEnv("JWT_SECRETS_MAP")
	currentKey := util.MustGetEnv("JWT_SECRETS_CURRENT_KEY")

	var secretsMap map[string]string
	err := json.Unmarshal([]byte(rawSecretsMap), &secretsMap)
	if err != nil {
		logrus.WithError(err).Fatal("Raw secrets map invalid JSON - not map[string]string")
	}

	jwtsecrets.Init(secretsMap, currentKey)
}
