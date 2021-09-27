package httpclient

type (
	nopObserver struct {
	}
)

func NewNopObserver() *nopObserver {
	return &nopObserver{}
}

func (n *nopObserver) Observe(f float64) {
}
