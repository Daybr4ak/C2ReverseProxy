package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"time"
)

func checkerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		tmp := make([]byte, 4096)
		dataErr := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if dataErr != nil {
			break
		}
		data, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println(data)
		fmt.Println(string(tmp[:data]))
		// fmt.Println(conn.Read(tmp))
		fmt.Println(reflect.TypeOf(data))
		// if data {
		// 	fmt.Println(1)
		// _, err = conn.Write([]byte("TO:SENDHello\r\n"))
		// if err != nil {
		// 	if err != io.EOF {
		// 		fmt.Println("read error:", err)
		// 	}
		// 	break
		// }
		// }
		// } else {
		// 	fmt.Println(1)
		// 	break
		// }
		// break
		_, err = conn.Write(tmp[:data])
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
	}
	fmt.Println(3)
}

func main() {
	Listen, err := net.Listen("tcp", "0.0.0.0:3997")
	checkerr(err)
	for {
		conn, err := Listen.Accept()
		checkerr(err)
		go handleConn(conn)
	}
}
