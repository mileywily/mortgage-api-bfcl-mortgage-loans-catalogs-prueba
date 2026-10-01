package http_fif

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
)

type RestClient interface {
	Get(ctx context.Context, path string, headers http.Header, params ...string) (*http.Response, error)
	Post(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error)
	Put(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error)
	Patch(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error)
	Delete(ctx context.Context, path string, headers http.Header) (*http.Response, error)
}

type DataEncoder interface {
	ContentType() string
	Encode(interface{}) (*bytes.Buffer, error)
}

type restClient struct {
	httpclient HTTPClient
	encoder    DataEncoder
}

func NewRestClient(opts ...interface{}) RestClient {
	var httpOpts []HttpClientOption
	var restOpts []RestClientOption

	for _, opt := range opts {
		switch o := opt.(type) {
		case HttpClientOption:
			httpOpts = append(httpOpts, o)
		case RestClientOption:
			restOpts = append(restOpts, o)
		default:
			panic(fmt.Errorf("opción desconocida"))
		}
	}

	rc := restClient{
		httpclient: NewHTTPClient(httpOpts...),
		encoder:    &JsonEncoder{},
	}

	for _, opt := range restOpts {
		rc = opt(rc)
	}
	return rc
}

// Get sends an HTTP GET request to the specified path with optional query parameters.
func (c restClient) Get(ctx context.Context, path string, headers http.Header, params ...string) (*http.Response, error) {
	req, err := c.makeRequest(http.MethodGet, path, nil, headers, params...)
	if err != nil {
		return nil, err
	}
	return c.httpclient.Do(req.WithContext(ctx))
}

// Post sends an HTTP POST request to the specified path with the provided body.
func (c restClient) Post(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	req, err := c.makeRequest(http.MethodPost, path, body, headers)
	if err != nil {
		return nil, err
	}
	return c.httpclient.Do(req.WithContext(ctx))
}

// Put sends an HTTP PUT request to the specified path with the provided body.
func (c restClient) Put(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	req, err := c.makeRequest(http.MethodPut, path, body, headers)
	if err != nil {
		return nil, err
	}
	return c.httpclient.Do(req.WithContext(ctx))
}

// Patch sends an HTTP PATCH request to the specified path with the provided body.
func (c restClient) Patch(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	req, err := c.makeRequest(http.MethodPatch, path, body, headers)
	if err != nil {
		return nil, err
	}
	return c.httpclient.Do(req.WithContext(ctx))
}

// Delete sends an HTTP DELETE request to the specified path.
func (c restClient) Delete(ctx context.Context, path string, headers http.Header) (*http.Response, error) {
	req, err := c.makeRequest(http.MethodDelete, path, nil, headers)
	if err != nil {
		return nil, err
	}
	return c.httpclient.Do(req.WithContext(ctx))
}

// makeRequest creates an HTTP request with a JSON body and optional query parameters.
func (c restClient) makeRequest(method, path string, body interface{}, headers http.Header, params ...string) (*http.Request, error) {
	encodedBody, err := c.encoder.Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, encodedBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add(contentTypeHeader, c.encoder.ContentType())

	for k, v := range headers {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}

	if len(params) != 0 {
		v, err := genQuery(params)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = v.Encode()
	}

	return req, nil
}
