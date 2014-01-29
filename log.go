package phosphorus

import "fmt"

var Logger *ServerLogger = &ServerLogger{}

type ServerLogger struct {
}

func ClientErrorf(c *Client, format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
