package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders_Parse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header
	headers = NewHeaders()
	data = []byte("       H@st : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("     HosT: localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 31, n)
	assert.False(t, done)

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("HOst: localhost:42069\r\nFriend:mauricio\r\n")

	// First header
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069", headers["host"])

	// advance buffer
	data = data[n:]

	// Second header
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "mauricio", headers["friend"])

	// Valid done
	headers = NewHeaders()
	data = []byte("\r\n")

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
	assert.Len(t, headers, 0)

}
