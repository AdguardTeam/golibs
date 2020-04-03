// Package jsonutil provides some helper functions for working with JSON objects
package jsonutil

import (
	ejson "encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

// JSON object
type JSON struct {
	m map[string]interface{}
}

// Exists returns TRUE if a key was specified in JSON source data
func (j *JSON) Exists(key string) bool {
	_, ok := j.m[key]
	return ok
}

// DecodeObject JSON data into a Go object and map-object
// A Go object has no method to check if the property was actually specified in JSON or not,
//  but a map-object provides this functionality.
// Note: not suitable for a large data
// obj: target object
// r: input data (reader object)
func DecodeObject(obj interface{}, r io.ReadCloser) (*JSON, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return DecodeObjectBuffer(obj, data)
}

// DecodeObjectBuffer - parse JSON data into a Go object and map-object
// obj: target object
// data: input data
func DecodeObjectBuffer(obj interface{}, data []byte) (*JSON, error) {
	reader := strings.NewReader(string(data))
	err := ejson.NewDecoder(reader).Decode(obj)
	if err != nil {
		return nil, err
	}

	return DecodeBuffer(data)
}

// DecodeBuffer - parse JSON data into a map-object
// data: input data
func DecodeBuffer(data []byte) (*JSON, error) {
	j := JSON{}
	j.m = make(map[string]interface{})
	err := ejson.Unmarshal(data, &j.m)
	if err != nil {
		return nil, err
	}
	return &j, nil
}
