package server

import (
	"log"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{log.New(log.Writer(), "", 0)}
}

func (l *Logger) Info(status int, path string, start time.Time) {
	l.Printf("\033[42m %d \033[0m | PATH: \033[33m\"%s\"\033[0m | DURATION: \033[42m %v \033[0m",
		status, path, time.Since(start),
	)
}

func (l *Logger) Error(status int, path string, err error) {
	l.Printf("\033[41m %d \033[0m | PATH: \033[33m\"%s\"\033[0m | ERROR: \033[31m %v \033[0m",
		status, path, err,
	)
}
