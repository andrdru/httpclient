package httpclient

type (
	Observer interface {
		Observe(float64)
	}
	nopObserver struct {
	}
)

func NewNopObserver() *nopObserver {
	return &nopObserver{}
}

func (n *nopObserver) Observe(f float64) {
}
