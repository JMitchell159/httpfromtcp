package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple headers
	headers = NewHeaders()
	data = []byte("HOST: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	data = []byte("hOsT: localhost:8080\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:8080", headers["host"])
	assert.Equal(t, 22, n)
	assert.False(t, done)
	data = []byte("HosT: localhost:8000\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, localhost:8080, localhost:8000", headers["host"])
	assert.Equal(t, 22, n)
	assert.False(t, done)
}
