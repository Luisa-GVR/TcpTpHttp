package response

import (
	"errors"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	Code200 StatusCode = iota
	Code400
	Code500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	str := ""

	switch statusCode {
	case Code200:
		str = "HTTP/1.1 200 OK\r\n"
	case Code400:
		str = "HTTP/1.1 400 Bad Request\r\n"
	case Code500:
		str = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		str = ""
	}

	_, err := w.Write([]byte(str))

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)

	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	contentLength, ok := headers.Get("Content-Length")
	if !ok {
		return errors.New("missing Content-Length")
	}

	Connection, ok := headers.Get("Connection")
	if !ok {
		return errors.New("missing Connection")
	}

	ContentType, ok := headers.Get("Content-Type")
	if !ok {
		return errors.New("missing Content-Type")
	}

	fullStr := "Content-Length: " + contentLength + "\r\nConnection: " +
		Connection + "\r\nContent-Type: " + ContentType + "\r\n\r\n"

	_, err := w.Write([]byte(fullStr))
	if err != nil {
		return err
	}
	return nil
}
