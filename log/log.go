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
	fullMessage := fmt.Sprintf("%s %s (%s) %s", time.Now().Format("2006-01-02 15:04:05.000 -0700 MST"), level, id, formattedInput)

	l.logger.Println(fullMessage)
}

func (l *Logger) ClientLog(level string, c interfaces.Client, id, format string, args ...interface{}) {
    var clientPrefix string
    if c.Account() != nil {
		clientPrefix = fmt.Sprintf("[%s / %s]", c.Connection().RemoteAddr().String(), c.Account().Name())
	} else {
		clientPrefix = fmt.Sprintf("[%s]", c.Connection().RemoteAddr().String())
	}

	formattedInput := fmt.Sprintf(format, args...)
	fullMessage := fmt.Sprintf("%s %s %s (%s) %s", time.Now().Format("2006-01-02 15:04:05.000 -0700 MST"), level, clientPrefix, id, formattedInput)

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

func (l *Logger) ClientDebug(c interfaces.Client, id, format string, args ...interface{}) {
	l.ClientLog("DEBUG", c, id, format, args...)
}

func (l *Logger) ClientInfo(c interfaces.Client, id, format string, args ...interface{}) {
	l.ClientLog("INFO ", c, id, format, args...)
}

func (l *Logger) ClientWarn(c interfaces.Client, id, format string, args ...interface{}) {
	l.ClientLog("WARN ", c, id, format, args...)
}

func (l *Logger) ClientError(c interfaces.Client, id, format string, args ...interface{}) {
	l.ClientLog("ERROR", c, id, format, args...)
}
