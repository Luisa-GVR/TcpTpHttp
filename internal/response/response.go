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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerStatus != stateInit {
		return errors.New("status line already written")
	}

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

	_, err := w.out.Write([]byte(str))

	if err != nil {
		return err
	}

	w.status = statusCode
	w.writerStatus = stateStatusWritten

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)

	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.writerStatus != stateHeadersWritten {
		return 0, errors.New("headers must be written before body")
	}

	n, err := w.out.Write(p)
	if err != nil {
		return 0, err
	}

	w.writerStatus = stateBodyWritten

	return n, nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.writerStatus != stateStatusWritten {
		return errors.New("status line must be written before headers")
	}

	for key, value := range headers {

		line := key + ":" + value + "\r\n"
		if _, err := w.out.Write([]byte(line)); err != nil {
			return err
		}
	}

	if _, err := w.out.Write([]byte("\r\n")); err != nil {
		return err
	}

	w.writerStatus = stateHeadersWritten
	return nil
}

type WriterStatus int

const (
	stateInit WriterStatus = iota
	stateStatusWritten
	stateHeadersWritten
	stateBodyWritten
)

type Writer struct {
	writerStatus WriterStatus
	out          io.Writer
	status       StatusCode
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{
		out:          out,
		writerStatus: stateInit,
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	if w.writerStatus != stateHeadersWritten && w.writerStatus != stateBodyWritten {
		return 0, errors.New("status line must be written before body")
	}

	sizeLineHex := strconv.FormatInt(int64(len(p)), 16) + "\r\n"

	_, err := w.out.Write([]byte(sizeLineHex))
	if err != nil {
		return 0, err
	}

	write, err := w.out.Write(p)
	if err != nil {
		return 0, err
	}

	_, err = w.out.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	w.writerStatus = stateBodyWritten

	return write, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	if w.writerStatus != stateHeadersWritten {
		return 0, errors.New("status line must be written before body")
	}

	w.writerStatus = stateBodyWritten
	write, err := w.out.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}

	return write, nil
}
