// logger/logger.go
package logger

import (
	"github.com/sirupsen/logrus"
	"time"
)

// LogUserAction logs a structured message about a user action
func LogUserAction(log *logrus.Logger, action string, userID string, page string, additionalFields map[string]interface{}) {
	// Create a structured log entry
	entry := log.WithFields(logrus.Fields{
		"timestamp": time.Now().Format(time.RFC3339),
		"user_id":   userID,
		"action":    action,
		"page":      page,
	})

	// Add any additional fields if provided
	if additionalFields != nil {
		for key, value := range additionalFields {
			entry = entry.WithField(key, value)
		}
	}

	// Log the action
	entry.Info("User action logged")
}
