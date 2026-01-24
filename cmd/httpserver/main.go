package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {

	handler := func(w *response.Writer, req *request.Request) {

		switch {

		case strings.HasPrefix(req.RequestLine.RequestTarget, "/yourproblem"):
			body := []byte(`<html>
							  <head>
								<title>400 Bad Request</title>
							  </head>
							  <body>
								<h1>Bad Request</h1>
								<p>Your request honestly kinda sucked.</p>
							  </body>
							</html>`)

			headers := response.GetDefaultHeaders(len(body))
			headers.Set("content-type", "text/html")

			w.WriteStatusLine(response.Code400)
			w.WriteHeaders(headers)
			w.WriteBody(body)

		case strings.HasPrefix(req.RequestLine.RequestTarget, "/myproblem"):
			body := []byte(`<html>
							  <head>
								<title>500 Internal Server Error</title>
							  </head>
							  <body>
								<h1>Internal Server Error</h1>
								<p>Okay, you know what? This one is on me.</p>
							  </body>
							</html>`)

			headers := response.GetDefaultHeaders(len(body))
			headers.Set("content-type", "text/html")

			w.WriteStatusLine(response.Code500)
			w.WriteHeaders(headers)
			w.WriteBody(body)

		case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin"):
			path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
			link := "https://httpbin.org" + path

			println(link)

			resp, err := http.Get(link)
			if err != nil {
				w.WriteStatusLine(response.Code500)
				return
			}
			defer resp.Body.Close()

			headers := response.GetDefaultHeaders(0)
			headers.Del("content-length")
			headers.Set("transfer-encoding", "chunked")

			w.WriteStatusLine(response.Code200)
			w.WriteHeaders(headers)

			buf := make([]byte, 1024)

			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					w.WriteChunkedBody(buf[:n])
					log.Println("forwarding", n, "bytes")
				}

				if err == io.EOF {
					break
				}
				if err != nil {
					return
				}
			}

			w.WriteChunkedBodyDone()
			return

		default:
			body := []byte(`<html>
							  <head>
								<title>200 OK</title>
							  </head>
							  <body>
								<h1>Success!</h1>
								<p>Your request was an absolute banger.</p>
							  </body>
							</html>`)

			headers := response.GetDefaultHeaders(len(body))
			headers.Set("content-type", "text/html")

			w.WriteStatusLine(response.Code200)
			w.WriteHeaders(headers)
			w.WriteBody(body)
		}
	}

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
