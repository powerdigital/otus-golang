package logger

import (
	"testing"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		config := config.LoggerConf{
			Level: levelInfo,
			File:  "/tmp/test.log",
		}

		logger := New(config)

		require.Equal(t, logger.level, config.Level)
		require.Equal(t, logger.logFile, config.File)
	})
}
