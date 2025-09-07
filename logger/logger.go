package logger

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{log.New(log.Writer(), "neural-network: ", 0)}
}

func (l *Logger) RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("\033[44m %s \033[0m | PATH: \033[33m\"%s\"\033[0m | TIME: %v",
			r.Method, r.URL.Path, start.Format("2006-01-02 15:04:05"))
		next.ServeHTTP(w, r)
	})
}

func (l *Logger) WithField(key string, value interface{}) {
	l.Printf("\033[42m INFO \033[0m | %s: %v", strings.ToUpper(key), value)
}

func (l *Logger) Debug(funcName string, obj interface{}) {
	l.Printf("\033[42m DEBUG \033[0m | FUNC: \033[33m\"%s\"\033[0m | OBJECT: \033[32m %v \033[0m",
		funcName, obj,
	)
}

func (l *Logger) Info(status int, path string, start time.Time) {
	l.Printf("\033[42m INFO \033[0m | STATUS: \033[32m%d\033[0m | PATH: \033[33m\"%s\"\033[0m | DURATION: \033[32m%s\033[0m",
		status, path, time.Since(start),
	)
}

func (l *Logger) Error(status int, path string, err error) {
	l.Printf("\033[41m ERROR \033[0m | STATUS: \033[31m%d\033[0m | PATH: \033[33m\"%s\"\033[0m | ERROR: \033[31m%s\033[0m",
		status, path, err,
	)
}

func (l *Logger) Fatal(err error) {
	l.Fatalf("\033[41m FATAL \033[0m | ERROR: \033[31m%s\033[0m", err)
}
