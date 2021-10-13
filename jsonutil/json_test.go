package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type jsonStruct struct {
	Str  string `json:"keyStr"`
	Bool bool   `json:"keyBool"`
}

func TestDecode(t *testing.T) {
	t.Parallel()

	o := jsonStruct{}
	data := `{"keyStr": "value", "keyBool": true}`
	j, err := DecodeObjectBuffer(&o, []byte(data))
	assert.True(t, err == nil)
	assert.True(t, o.Str == "value")
	assert.True(t, o.Bool == true)
	assert.True(t, j.Exists("keyStr"))
	assert.True(t, j.Exists("keyBool"))

	o = jsonStruct{}
	data = `{}`
	j, err = DecodeObjectBuffer(&o, []byte(data))
	assert.True(t, err == nil)
	assert.True(t, o.Str == "")
	assert.True(t, o.Bool == false)
	assert.True(t, !j.Exists("keyStr"))
	assert.True(t, !j.Exists("keyBool"))
}
