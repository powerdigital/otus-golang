package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/powerdigital/otus-golang/hw12_13_14_15_calendar/internal/config"
)

const (
	levelInfo     = "INFO"
	levelError    = "ERROR"
	levelCritical = "CRITICAL"
	levelFatal    = "FATAL"
)

type Logger struct {
	level   string
	logFile string
}

type LogHandler interface {
	Info(msg string)
	Error(msg string)
	Critical(msg string)
	Fatal(msg string)
}

func New(config config.LoggerConf) *Logger {
	return &Logger{
		level:   config.Level,
		logFile: config.File,
	}
}

func (l Logger) Info(msg string) {
	log := fmt.Sprintf("[%s] %s\n", levelInfo, msg)
	l.WriteLog(log)
}

func (l Logger) Error(msg string) {
	log := fmt.Sprintf("[%s] %s\n", levelError, msg)
	l.WriteLog(log)
}

func (l Logger) Critical(msg string) {
	log := fmt.Sprintf("[%s] %s\n", levelCritical, msg)
	l.WriteLog(log)
}

func (l Logger) Fatal(msg string) {
	log := fmt.Sprintf("[%s] %s\n", levelFatal, msg)
	l.WriteLog(log)
}

func (l Logger) WriteLog(msg string) {
	timeCurr := time.Now()
	msg = strings.Join([]string{timeCurr.Format("2006-01-02 15:04:05"), msg}, " ")

	f, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(msg); err != nil {
		panic(err)
	}
}
