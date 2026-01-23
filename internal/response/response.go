package response

import (
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
		str = "HTTP/1.1 200 OK\r\n"
	}

	_, err := w.Write([]byte(str))

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)

	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	contentLength, _ := headers.Get("content-length")
	connection, _ := headers.Get("connection")
	contentType, _ := headers.Get("content-type")

	if contentLength == "" {
		contentLength = "0"
	}
	if connection == "" {
		connection = "close"
	}
	if contentType == "" {
		contentType = "text/plain"
	}

	fullStr := "content-length: " + contentLength + "\r\n" +
		"connection: " + connection + "\r\n" +
		"content-type: " + contentType + "\r\n\r\n"

	_, err := w.Write([]byte(fullStr))
	return err
}
