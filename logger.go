package httpclient

type (
	LoggerLevel int

	Logger interface {
		Printf(format string, v ...interface{})
	}

	nopLogger struct {
	}
)

const (
	LoggerLevelNop   LoggerLevel = 0
	LoggerLevelError LoggerLevel = 1
	LoggerLevelInfo  LoggerLevel = 2
)

func NewNopLogger() *nopLogger {
	return &nopLogger{}
}

func (n *nopLogger) Printf(format string, v ...interface{}) {
}

func (l LoggerLevel) Allowed(level LoggerLevel) bool {
	return level >= l
}
