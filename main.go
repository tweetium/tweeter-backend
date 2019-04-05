package main

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver
	"github.com/sirupsen/logrus"

	"tweeter/db"
	"tweeter/handlers"
	"tweeter/util"
)

func main() {
	port := util.MustGetEnvUInt32("PORT")
	dbURL := util.MustGetEnv("DATABASE_URL")

	sentry_dsn := os.Getenv("SENTRY_DSN")
	hook, err := logrus_sentry.NewSentryHook(sentry_dsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	hook.StacktraceConfiguration.Enable = true

	if err != nil {
		logrus.WithFields(logrus.Fields{"dsn": sentry_dsn}).Error("Failed to set up Sentry hook")
	} else {
		logrus.WithFields(logrus.Fields{"dsn": sentry_dsn}).Info("Successfully set up Sentry hook")
		logrus.AddHook(hook)
	}

	err = db.Init(dbURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to initialize database")
	}

	handlers.RunWebserver(port)
}
