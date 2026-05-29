// examples/internal/sampleevents/events.go
package sampleevents

import (
	"time"

	"github.com/smegg99/s99logger"
)

func AppStarted(service, version string) s99logger.Event {
	return s99logger.NewEvent("app_started",
		s99logger.String("service", service),
		s99logger.String("version", version),
	)
}

func LoginFailed(username string, err error) s99logger.Event {
	return s99logger.NewEvent("login_failed",
		s99logger.String("username", username),
		s99logger.Err(err),
	)
}

func ReconnectAttempt(host string, attempt, max int, backoff time.Duration) s99logger.Event {
	return s99logger.NewEvent("reconnect_attempt",
		s99logger.String("host", host),
		s99logger.Int("attempt", attempt),
		s99logger.Int("max", max),
		s99logger.Duration("backoff", backoff),
	)
}

func JobsProcessed(count int) s99logger.Event {
	return s99logger.NewPluralEvent("jobs_processed", count, s99logger.Int("count", count))
}
