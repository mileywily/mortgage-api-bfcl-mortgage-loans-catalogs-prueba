package http_fif

import (
	"bytes"
	"io"
	"net/url"
	"strings"
)

type FormURLEncoder struct {
}

func (e *FormURLEncoder) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (e *FormURLEncoder) Encode(body interface{}) (*bytes.Buffer, error) {
	data, ok := body.(map[string]string)
	if !ok {
		return nil, ErrInvalidBody
	}
	var formBody string
	for key, value := range data {
		if formBody != "" {
			formBody += "&"
		}
		formBody += url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, strings.NewReader(formBody)); err != nil {
		return nil, ErrInvalidBody
	}

	return buf, nil
}
