package logging

import (
	"testing"

	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		level       string
		expectedFmt logrus.Formatter
		expectedLvl logrus.Level
	}{
		{
			name:        "test json format and debug level",
			format:      "json",
			level:       "debug",
			expectedFmt: &logrus.JSONFormatter{},
			expectedLvl: logrus.DebugLevel,
		},
		{
			name:        "test json format and warning level",
			format:      "json",
			level:       "warning",
			expectedFmt: &logrus.JSONFormatter{},
			expectedLvl: logrus.WarnLevel,
		},
		{
			name:        "test text format and error level",
			format:      "text",
			level:       "error",
			expectedFmt: &logrus.TextFormatter{},
			expectedLvl: logrus.ErrorLevel,
		},
		{
			name:        "test text format and info level",
			format:      "text",
			level:       "info",
			expectedFmt: &logrus.TextFormatter{},
			expectedLvl: logrus.InfoLevel,
		},
		{
			name:        "test unknown format and level",
			format:      "unknown",
			level:       "unknown",
			expectedFmt: &logrus.TextFormatter{},
			expectedLvl: logrus.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appConfig := &config.AppConfig{
				Logging: config.Logging{
					Format: tt.format,
					Level:  tt.level,
				},
			}

			logger := NewLogger(appConfig)

			assert.Equal(t, tt.expectedFmt, logger.Formatter)
			assert.Equal(t, tt.expectedLvl, logger.Level)
		})
	}
}
