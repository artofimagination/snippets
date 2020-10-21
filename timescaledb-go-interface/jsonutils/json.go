package jsonutils

import (
	"bytes"
	"encoding/json"
	"io"
)

// ProcessJSON returns the parsed json content.
func ProcessJSON(data []byte) (map[string]interface{}, error) {
	// use it
	d := make(map[string]interface{})
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}
