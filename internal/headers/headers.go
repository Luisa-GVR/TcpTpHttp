package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	fullValue := string(data)

	//not enough data to parse
	if !strings.Contains(fullValue, "\r\n") {
		return 0, false, nil
	}
	//CRLF at the start, end of headers
	if strings.HasPrefix(fullValue, "\r\n") {
		return 2, true, nil
	}

	beforeCRLFIndex := strings.Index(fullValue, "\r\n")
	beforeCRLF := fullValue[:beforeCRLFIndex]

	before, after, found := strings.Cut(beforeCRLF, ":")
	if !found {
		return 0, false, errors.New("invalid header line")
	}

	// Only whitespace allowed is the suffix/leading whitespace before the field name
	key := strings.TrimLeft(before, " \t")

	// Any other type of white space is invalid
	if strings.ContainsAny(key, " \t") {
		return 0, false, errors.New("invalid header key")
	}
	// lowercase everything
	key = strings.ToLower(key)

	/*
		Last check on the key, must be from a-z (already lower-cased it), 0-9 and the following special characters
		!, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
	*/

	allowedCharsString := "abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	for _, char := range key {
		if !strings.ContainsRune(allowedCharsString, char) {
			return 0, false, errors.New("invalid header line")
		}
	}

	after = strings.TrimSpace(after)

	bytesConsumed := len(beforeCRLF) + 2

	//Lets check if the header already exist
	if _, exist := h[key]; exist {
		h[key] += "," + after
	} else {
		h[key] = after

	}

	//Remember: The Parse function should only return done=true when the data starts with a CRLF, which can't happen when it finds a new key/value pair.

	return bytesConsumed, false, nil
}

func (h Headers) Get(key string) (value string, ok bool) {

	key = strings.ToLower(key)
	value, ok = h[key]
	return value, ok
}

func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	h[key] = value

}
