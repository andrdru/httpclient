package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
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
		logLevel  LoggerLevel
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
		log:         NewNopLogger(),
		loggerLevel: LoggerLevelNop,
		timeout:     TimeoutDefault,
		scheme:      SchemeDefault,
		rateLimit:   NewNopRateLimit(),
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
		log:       args.log,
		logLevel:  args.loggerLevel,
		rateLimit: args.rateLimit,
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
	c.infoLog("http request: %s\n", req)

	var rs *http.Response
	rs, err = c.client.Do(req)
	if err != nil {
		c.errorLog("http request error: %s\n", err.Error())
		return 0, []byte{}, err
	}

	defer func() {
		if err = rs.Body.Close(); err != nil {
			c.errorLog("http body close error %s\n", wrapErr(err, closeError).Error())
		}
	}()

	statusCode = rs.StatusCode
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		c.errorLog("http body read error: %s\n", err.Error())
		return 0, []byte{}, err
	}
	c.infoLog("http response: %s\n", string(body))

	return statusCode, body, nil
}

func (c *httpClient) errorLog(format string, values ...interface{}) {
	if c.logLevel.Allowed(LoggerLevelError) {
		c.log.Printf(format, values...)
	}
}

func (c *httpClient) infoLog(format string, values ...interface{}) {
	if !c.logLevel.Allowed(LoggerLevelInfo) {
		return
	}

	for i := range values {
		var result = logMarshalled(values[i])
		if result != "" {
			values[i] = result
		}
	}
	c.log.Printf(format, values...)
}

func logMarshalled(v interface{}) string {
	var b, err = json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
