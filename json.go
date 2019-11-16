package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	JsonClient interface {
		Delete(
			ctx context.Context,
			path string,
			param map[string]string,
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Get(
			ctx context.Context,
			path string,
			param map[string]string,
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Patch(
			ctx context.Context,
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Post(
			ctx context.Context,
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
		Put(
			ctx context.Context,
			path string,
			param map[string]interface{},
			header http.Header,
		) (statusCode int, res interface{}, err error)
	}

	jsonClient struct {
		HttpClient
	}
)

const (
	headerContentType        = "Content-Type"
	headerContentTypeDefault = "application/json"
)

func NewJsonClient(host string, options ...Option) *jsonClient {
	return &jsonClient{
		HttpClient: NewHttpClient(host, options...),
	}
}

func (c *jsonClient) Delete(
	ctx context.Context,
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.requestNoBody(ctx, http.MethodDelete, path, param, header)
}

func (c *jsonClient) Get(
	ctx context.Context,
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.requestNoBody(ctx, http.MethodGet, path, param, header)
}

func (c *jsonClient) Patch(
	ctx context.Context,
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(ctx, http.MethodPatch, path, param, header)
}

func (c *jsonClient) Post(
	ctx context.Context,
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(ctx, http.MethodPost, path, param, header)
}

func (c *jsonClient) Put(
	ctx context.Context,
	path string,
	param map[string]interface{},
	header http.Header,
) (statusCode int, res interface{}, err error) {
	return c.request(ctx, http.MethodPut, path, param, header)
}

func (c *jsonClient) request(
	ctx context.Context,
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

	if header.Get(headerContentType) == "" {
		header.Set(headerContentType, headerContentTypeDefault)
	}

	statusCode, body, err = c.HttpClient.Request(ctx, method, path, body, header)
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
	ctx context.Context,
	method string,
	path string,
	param map[string]string,
	header http.Header,
) (statusCode int, res interface{}, err error) {
	var body []byte

	if header.Get(headerContentType) == "" {
		header.Set(headerContentType, headerContentTypeDefault)
	}

	statusCode, body, err = c.HttpClient.RequestNoBody(ctx, method, path, param, header)
	if err != nil {
		return 0, nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, nil, err
	}

	return statusCode, res, nil
}
