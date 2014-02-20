package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var baseLogger *log.Logger = log.New(os.Stdout, "", 0)
var Server *Logger = &Logger{
	logger: baseLogger,
}

type Logger struct {
	logger   *log.Logger
	prefixer func(string) string
}

func NewLogger() *Logger {
	return &Logger{logger: baseLogger}
}

func (l *Logger) SetPrefixer(f func(string) string) {
	l.prefixer = f
}

func (l *Logger) Log(level, id, format string, args ...interface{}) {
	formattedInput := fmt.Sprintf(format, args...)
	baseMessage := fmt.Sprintf("(%s) %s", id, formattedInput)

	if l.prefixer != nil {
		baseMessage = l.prefixer(baseMessage)
	}

	fullMessage := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05.000 -0700 MST"), level, baseMessage)

	l.logger.Println(fullMessage)
}

func (l *Logger) Debug(id, format string, args ...interface{}) {
	l.Log("DEBUG", id, format, args...)
}

func (l *Logger) Info(id, format string, args ...interface{}) {
	l.Log("INFO ", id, format, args...)
}

func (l *Logger) Warn(id, format string, args ...interface{}) {
	l.Log("WARN ", id, format, args...)
}

func (l *Logger) Error(id, format string, args ...interface{}) {
	l.Log("ERROR", id, format, args...)
}
