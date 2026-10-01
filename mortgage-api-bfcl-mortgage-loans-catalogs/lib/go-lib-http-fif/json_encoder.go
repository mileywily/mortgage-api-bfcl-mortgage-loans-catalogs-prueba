package http_fif

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// encodeBody encodes the provided interface into a JSON byte buffer.
// It validates the structure before encoding and returns an error if the validation or encoding fails.
type JsonEncoder struct {
}

func (e *JsonEncoder) ContentType() string {
	return "application/json"
}

func (e *JsonEncoder) Encode(body interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if body == nil {
		return buf, nil
	}

	err := validate.Struct(body)
	if err != nil {
		return nil, fmt.Errorf("can't encode %v: %w", err, ErrFailedValidation)
	}

	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		return buf, fmt.Errorf("can't encode %v: %w", body, err)
	}
	return buf, nil
}
