// Package jsonutil provides some helper functions for working with JSON.
package jsonutil

import (
	"encoding/json"
	"io"
	"strings"
)

// JSON object
//
// Deprecated: This type is deprecated and will be removed in a future release.
type JSON struct {
	m map[string]any
}

// Exists returns TRUE if a key was specified in JSON source data
func (j *JSON) Exists(key string) bool {
	_, ok := j.m[key]
	return ok
}

// DecodeObject JSON data into a Go object and map-object
// A Go object has no method to check if the property was actually specified in
// JSON or not, but a map-object provides this functionality.
// Note: not suitable for a large data
// obj: target object
// r: input data (reader object)
func DecodeObject(obj any, r io.ReadCloser) (*JSON, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return DecodeObjectBuffer(obj, data)
}

// DecodeObjectBuffer - parse JSON data into a Go object and map-object
// obj: target object
// data: input data
func DecodeObjectBuffer(obj any, data []byte) (*JSON, error) {
	reader := strings.NewReader(string(data))
	err := json.NewDecoder(reader).Decode(obj)
	if err != nil {
		return nil, err
	}

	return DecodeBuffer(data)
}

// DecodeBuffer - parse JSON data into a map-object
// data: input data
func DecodeBuffer(data []byte) (*JSON, error) {
	j := JSON{}
	j.m = make(map[string]any)
	err := json.Unmarshal(data, &j.m)
	if err != nil {
		return nil, err
	}

	return &j, nil
}
