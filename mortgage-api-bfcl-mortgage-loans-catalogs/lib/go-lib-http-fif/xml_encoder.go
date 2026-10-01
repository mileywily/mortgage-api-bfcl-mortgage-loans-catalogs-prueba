package http_fif

import (
	"bytes"
	"encoding/xml"
)

type XmlEncoder struct{}

func (e *XmlEncoder) ContentType() string {
	return "application/xml"
}

func (e *XmlEncoder) Encode(v interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}
