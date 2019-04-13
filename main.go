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

	sentryDSN := os.Getenv("SENTRY_DSN")
	hook, err := logrus_sentry.NewSentryHook(sentryDSN, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.Level = logrus.WarnLevel

	if err != nil {
		logrus.WithField("dsn", sentryDSN).Error("Failed to set up Sentry hook")
	} else {
		logrus.WithField("dsn", sentryDSN).Info("Successfully set up Sentry hook")
		logrus.AddHook(hook)
	}

	err = db.Init(dbURL)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}

	handlers.RunWebserver(port)
}
