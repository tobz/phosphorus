package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tobz/phosphorus/interfaces"
)

var Server *Logger = &Logger{
	logger: log.New(os.Stdout, "", 0),
}

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Log(level, id, format string, args ...interface{}) {
	formattedInput := fmt.Sprintf(format, args...)
	fullMessage := fmt.Sprintf("%s [%s] (%s) %s", time.Now().Format("2006-01-02 15:04:05.000 -0700 MST"), level, id, formattedInput)

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

func WrapForClient(c interfaces.Client, format string, args ...interface{}) string {
	var prefix string

	if c.Account() != nil {
		prefix = fmt.Sprintf("[%s / %s] ", c.RemoteAddr().String(), c.Account().Name())
	} else {
		prefix = fmt.Sprintf("[%s] ", c.RemoteAddr().String())
	}

	return fmt.Sprintf(prefix+format, args...)
}
