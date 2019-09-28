package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	HttpClient interface {
		RequestNoBody(method string, path string, param *map[string]string, header *map[string]string) (statusCode int, body []byte, err error)
		Request(method string, path string, param []byte, header *map[string]string) (statusCode int, body []byte, err error)
	}

	httpClient struct {
		client http.Client
		scheme string
		host   string
	}
)

func NewHttpClient(scheme string, host string, timeout time.Duration) *httpClient {
	return &httpClient{
		scheme: scheme,
		host:   host,
		client: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *httpClient) RequestNoBody(method string, path string, param *map[string]string, header *map[string]string) (statusCode int, body []byte, err error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	if param != nil {
		q := u.Query()
		for key, val := range *param {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return 0, []byte{}, err
	}

	return c.request(req, header)
}

func (c *httpClient) Request(method string, path string, param []byte, header *map[string]string) (statusCode int, body []byte, err error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	r := bytes.NewReader(param)
	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return 0, []byte{}, err
	}

	return c.request(req, header)
}

func (c *httpClient) request(req *http.Request, header *map[string]string) (statusCode int, body []byte, err error) {
	if header != nil {
		for key, val := range *header {
			req.Header.Set(key, val)
		}
	}

	rs, err := c.client.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}

	defer func() {
		_ = rs.Body.Close()
	}()

	statusCode = rs.StatusCode
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return statusCode, body, nil
}
