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

	//Lets ensure theres no space between colon and key
	if strings.HasSuffix(before, " ") {
		return 0, false, errors.New("invalid header")
	}
	//Now lets ensure theres no whitespace remaining before checking everything
	before = strings.TrimSpace(before)
	after = strings.TrimSpace(after)

	bytesConsumed := len(beforeCRLF) + 2

	h[before] = after
	//Remember: The Parse function should only return done=true when the data starts with a CRLF, which can't happen when it finds a new key/value pair.

	return bytesConsumed, false, nil
}
