package http_fif

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPClientMiddleware func(HTTPClient) HTTPClient

// httpClient represents an HTTP requester with additional features such as logging and tracing.
type httpClient struct {
	baseClient *http.Client
	requester  HTTPClient
	headers    http.Header
	baseURL    *url.URL
	logger     HTTPClientMiddleware
	tracer     HTTPClientMiddleware
}

// NewHTTPClient creates a new httpClient with the provided options.
// If no options are provided, it uses default options.
func NewHTTPClient(options ...HttpClientOption) HTTPClient {
	cli := &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
	c := httpClient{
		baseClient: cli,
		requester:  cli,
		headers:    http.Header{},
		baseURL:    nil,
		logger:     nil,
		tracer:     nil,
	}

	var defaultOptions = []HttpClientOption{
		BaseURL(""),
		Timeout(15 * time.Second),
	}

	options = append(defaultOptions, options...)

	for _, option := range options {
		c = option(c)
	}

	if c.tracer != nil {
		c.requester = c.tracer(c.requester)
	}

	if c.logger != nil {
		c.requester = c.logger(c.requester)
	}

	return c
}

// Do sends an HTTP request and returns an HTTP response, handling logging and tracing.
func (c httpClient) Do(req *http.Request) (res *http.Response, err error) {
	if req == nil {
		return nil, ErrRequestNil
	}

	url, err := resolveURL(c.baseURL, req.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrInvalidURLAddress)
	}
	req.URL = url

	for k, v := range c.headers {
		if old := req.Header.Get(k); old == "" {
			for _, s := range v {
				req.Header.Add(k, s)
			}
		}
	}

	//TODO: add request retries
	res, err = c.requester.Do(req)

	return res, err

}

// resolveURL resolves the final URL by combining the base URL and the path.
func resolveURL(base, path *url.URL) (*url.URL, error) {
	b := ""
	if base != nil {
		b = base.String()
	}

	p := ""
	if path != nil {
		p = path.String()
	}
	return url.Parse(b + p)
}

// genQuery generates URL query parameters from a list of key-value pairs.
func genQuery(params []string) (url.Values, error) {
	v := url.Values{}
	size := len(params)

	if size%2 != 0 {
		return nil, ErrInvalidParametersCount
	}

	for i := 0; i < size; i += 2 {
		v.Add(params[i], params[i+1])
	}
	return v, nil
}
