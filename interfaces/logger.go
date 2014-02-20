package interfaces

type Logger interface {
	Debug(string, string, ...interface{})
	Info(string, string, ...interface{})
	Warn(string, string, ...interface{})
	Error(string, string, ...interface{})
}
