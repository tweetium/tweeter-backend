package log

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// Init sets up the logrus sentry hook
func Init() {
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
}
