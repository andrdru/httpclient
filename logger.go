package httpclient

type (
	Logger interface {
		Printf(format string, v ...interface{})
	}

	nopLogger struct {
	}
)

func NewNopLogger() *nopLogger {
	return &nopLogger{}
}

func (n *nopLogger) Printf(format string, v ...interface{}) {
}
