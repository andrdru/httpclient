package httpclient

type (
	Logger interface {
		Println(v ...interface{})
	}

	nopLogger struct {
	}
)

func NewNopLogger() *nopLogger {
	return &nopLogger{}
}

func (n *nopLogger) Println(v ...interface{}) {
}
