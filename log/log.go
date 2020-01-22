package log

import (
	"fmt"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// Init sets up the logrus sentry hook
func Init(release string) {
	hook, err := logrus_sentry.NewAsyncWithClientSentryHook(raven.DefaultClient, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.Level = logrus.WarnLevel
	hook.SetRelease(release)

	if err != nil {
		logrus.WithError(err).Error("Failed to set up Sentry hook")
	} else {
		logrus.Info("Successfully set up Sentry hook")
		logrus.AddHook(hook)
		logrus.AddHook(&flushSentryHook{hook})
	}
}

// Internal type that flushes the sentry hook for panic / fatals
// See the issue on logrus_sentry repo here: https://github.com/evalphobia/logrus_sentry/issues/46
type flushSentryHook struct {
	hook *logrus_sentry.SentryHook
}

func (f *flushSentryHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
}

func (f *flushSentryHook) Fire(*logrus.Entry) error {
	// sinc we're async flushing, lets make sure we do a final flush
	// before we exit the program
	timeout := time.After(10 * time.Second)
	flushed := make(chan interface{}, 1)
	go func() {
		f.hook.Flush()
		flushed <- struct{}{}
	}()
	select {
	case <-flushed:
		return nil
	case <-timeout:
		return fmt.Errorf("Timed out after 10s waiting to flush sentry hook")
	}
}
