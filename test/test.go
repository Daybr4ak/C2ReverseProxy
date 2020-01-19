package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func checkerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleConn(conn net.Conn) {
	for {
		tmp := make([]byte, 4096)
		data, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println(tmp[:data])
		_, err = conn.Write([]byte("\r\n"))
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

	}
}

func main() {
	Listen, err := net.Listen("tcp", "0.0.0.0:8888")
	checkerr(err)
	for {
		conn, err := Listen.Accept()
		checkerr(err)
		go handleConn(conn)
	}
}
