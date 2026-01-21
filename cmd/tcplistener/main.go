package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	f, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not open")
	}
	defer f.Close()

	for {
		conn, err := f.Accept()
		if err != nil {
			log.Fatalf("could not accept")
			continue
		}
		fmt.Println("Yay! connected!:)")
		toPrint, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("could not read")
		}
		fmt.Println("Request line:")
		fmt.Println("- Method: " + toPrint.RequestLine.Method)
		fmt.Println("- Target: " + toPrint.RequestLine.RequestTarget)
		fmt.Println("- Version: " + toPrint.RequestLine.HttpVersion)
		fmt.Println("Headers:")

		for key, value := range toPrint.Headers {
			fmt.Println("-  " + key + ": " + value)
		}

		//fmt.Printf("connection closed, oops!")

	}

}
