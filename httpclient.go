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
		client    http.Client
		scheme    string
		host      string
		log       Logger
		rateLimit RateLimiter
	}
)

const (
	closeError = "close fails"

	TimeoutDefault = 30 * time.Second
	SchemeDefault  = "https"
)

func NewHttpClient(host string, options ...Option) *httpClient {
	var args = &Options{
		log:       NewNopLogger(),
		timeout:   TimeoutDefault,
		scheme:    SchemeDefault,
		rateLimit: NewNopRateLimit(),
	}

	var opt Option
	for _, opt = range options {
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
	var u = url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	if param != nil {
		var q = u.Query()
		var key, val string
		for key, val = range param {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
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
	var u = url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	var r = bytes.NewReader(param)
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return 0, []byte{}, err
	}

	return c.request(req, header)
}

func (c *httpClient) request(req *http.Request, header http.Header) (statusCode int, body []byte, err error) {
	c.rateLimit.Take()

	if header != nil {
		req.Header = header
	}

	var rs *http.Response
	rs, err = c.client.Do(req)
	if rs != nil {
		defer func() {
			if err = rs.Body.Close(); err != nil {
				c.warnErr(wrapErr(err, closeError))
			}
		}()
	}

	if err != nil {
		return 0, []byte{}, err
	}

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
