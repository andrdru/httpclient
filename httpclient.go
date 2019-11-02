package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	HttpClient interface {
		RequestNoBody(
			ctx context.Context,
			method string,
			path string,
			param map[string]string,
			header http.Header,
		) (statusCode int, body []byte, err error)
		Request(
			ctx context.Context,
			method string,
			path string,
			param []byte,
			header http.Header,
		) (statusCode int, body []byte, err error)
	}

	httpClient struct {
		client http.Client
		scheme string
		host   string
		log    Logger
	}
)

const (
	closeError = "close fails"

	TimeoutDefault = 5 * time.Second
	SchemeDefault  = "https"
)

func NewHttpClient(host string, options ...Option) *httpClient {
	args := &Options{
		log:     NewNopLogger(),
		timeout: TimeoutDefault,
		scheme:  SchemeDefault,
	}

	for _, opt := range options {
		opt(args)
	}

	return &httpClient{
		host:   host,
		scheme: args.scheme,
		client: http.Client{
			Timeout: args.timeout,
		},
		log: args.log,
	}
}

func (c *httpClient) RequestNoBody(
	ctx context.Context,
	method string,
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, body []byte, err error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	if param != nil {
		q := u.Query()
		for key, val := range param {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return 0, []byte{}, err
	}

	return c.request(req, header)
}

func (c *httpClient) Request(
	ctx context.Context,
	method string,
	path string,
	param []byte,
	header http.Header,
) (statusCode int, body []byte, err error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	r := bytes.NewReader(param)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return 0, []byte{}, err
	}

	return c.request(req, header)
}

func (c *httpClient) request(req *http.Request, header http.Header) (statusCode int, body []byte, err error) {
	if header != nil {
		req.Header = header
	}

	rs, err := c.client.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}

	defer func() {
		if err := rs.Body.Close(); err != nil {
			c.warnErr(wrapErr(err, closeError))
		}
	}()

	statusCode = rs.StatusCode
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return statusCode, body, nil
}

func (c *httpClient) warnErr(err error) {
	c.log.Println(err.Error())
}
