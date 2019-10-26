package httpclient

import (
	"encoding/json"
	"net/http"
)

type (
	JsonClient interface {
		Delete(
			path string,
			param map[string]string,
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Get(
			path string,
			param map[string]string,
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Patch(
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Post(
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Put(
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
	}

	jsonClient struct {
		HttpClient
	}
)

func NewJsonClient(host string, options ...Option) *jsonClient {
	return &jsonClient{
		HttpClient: NewHttpClient(host, options...),
	}
}

func (c *jsonClient) Delete(
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.requestNoBody(http.MethodDelete, path, param, header)
}

func (c *jsonClient) Get(
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.requestNoBody(http.MethodGet, path, param, header)
}

func (c *jsonClient) Patch(
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(http.MethodPatch, path, param, header)
}

func (c *jsonClient) Post(
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(http.MethodPost, path, param, header)
}

func (c *jsonClient) Put(
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(http.MethodPut, path, param, header)
}

func (c *jsonClient) request(
	method string,
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	var body []byte
	if param != nil {
		body, err = json.Marshal(param)
		if err != nil {
			return 0, nil, err
		}
	}
	statusCode, body, err = c.HttpClient.Request(method, path, body, header)
	if err != nil {
		return 0, nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, nil, err
	}

	return statusCode, res, nil
}

func (c *jsonClient) requestNoBody(
	method string,
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	var body []byte
	statusCode, body, err = c.HttpClient.RequestNoBody(method, path, param, header)
	if err != nil {
		return 0, nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, nil, err
	}

	return statusCode, res, nil
}
