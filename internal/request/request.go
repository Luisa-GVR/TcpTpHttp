package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type State int

const (
	Initialized State = iota
	Done
	requestStateParsingHeaders
	requestStateParsingBody
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	const bufferSize = 8
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	var request Request
	request.State = Initialized

	for request.State != Done {
		//buffer full, we make a new one
		if readToIndex >= cap(buf) {

			newBuf := make([]byte, len(buf)*2, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.State != Done {
					//This error should be enough for a content length being too short
					return nil, fmt.Errorf("unexpected EOF before request was fully parsed")
				}
				break
			}

			return nil, err
		}

		readToIndex += n
		consumed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}

	}

	return &request, nil
}

// First part of the request
func parseRequestLine(line string) (requestLine RequestLine, bytes int, err error) {

	var reqLine RequestLine
	if !strings.Contains(line, "\r\n") {
		return reqLine, 0, nil
	}

	eachRequestLine := strings.Split(line, "\r\n")
	words := strings.Fields(eachRequestLine[0])

	if len(words) < 3 {
		return reqLine, 0, fmt.Errorf("RequestLine must have at least three words")
	}
	//Method = 0, request targed = 1, httpversion =3

	//Method must be fully uppercased
	if !(words[0] == strings.ToUpper(words[0])) {
		return reqLine, 0, fmt.Errorf("method must be fully uppercase")
	}
	reqLine.Method = words[0]

	//Http version must be 1.1
	if words[2] != "HTTP/1.1" {
		return reqLine, 0, fmt.Errorf("RequestLine must use HTTP/1.1")
	}
	reqLine.HttpVersion = "1.1"

	//request-target must start with "/"
	if !strings.HasPrefix(words[1], "/") {
		return reqLine, 0, fmt.Errorf("request target must start with '/'")
	}
	reqLine.RequestTarget = words[1]

	lineEnd := strings.Index(line, "\r\n")
	bytesConsumed := lineEnd + 2

	return reqLine, bytesConsumed, nil
}

// Parser for packets
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != Done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		//needs more data
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.State == Initialized {
		reqLine, bytesConsumed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesConsumed == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.State = requestStateParsingHeaders
		return bytesConsumed, nil
	} else if r.State == requestStateParsingHeaders {

		if r.Headers == nil {
			r.Headers = headers.NewHeaders()
		}

		n, done, err := r.Headers.Parse(data)

		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		if done {
			r.State = requestStateParsingBody
		}
		return n, nil

	} else if r.State == requestStateParsingBody {

		contentLengthStr, ok := r.Headers.Get("Content-Length")

		if !ok {
			//no content length, so nothing to parse
			r.State = Done
			return 0, nil
		}

		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length")
		}

		r.Body = data

		if len(r.Body) == contentLength {
			r.State = Done
			fmt.Printf("Consumed the entire length of data")

		}
		if len(r.Body) > contentLength {
			return 0, fmt.Errorf("body larger than Content-Length")
		}

	} else if r.State == Done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	} else {
		return 0, fmt.Errorf("error: unknown state")
	}

	return 0, nil

}
