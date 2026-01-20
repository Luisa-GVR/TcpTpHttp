package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp4", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	udpConn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		str, errReader := reader.ReadString('\n')
		if errReader != nil {
			log.Fatal(errReader)
		}
		udpConn.Write([]byte(str))
	}

}
