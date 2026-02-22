package utils

import (
	"log"
	"os"
)

type CustomLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

func NewLogger(serviceName string) *CustomLogger {
	return &CustomLogger{
		infoLogger:  log.New(os.Stdout, "INFO ["+serviceName+"]: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stdout, "ERROR ["+serviceName+"]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *CustomLogger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *CustomLogger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}
