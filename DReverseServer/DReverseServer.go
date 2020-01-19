package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var (
	dataBuf = make([]byte, 0, 4096)
)

// 抛出异常
func checkerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// 获取数据
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
		if string(tmp[:6]) == "TO:GET" {
			_, err = conn.Write(dataBuf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				break
			}
			dataBuf = append(dataBuf[:0])
		} else if string(tmp[:10]) == "TO:CONNECT" {
			_, err = conn.Write([]byte("CONNECT SUCCESS\n"))
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				break
			}
		} else if string(tmp[:7]) == "TO:SEND" {
			fmt.Println("Get tmp:", tmp[7:data])
		} else {
			dataBuf = append(dataBuf, tmp[:data]...)
		}
		fmt.Println("tmp:", tmp[:data])
		fmt.Println("databuf:", dataBuf)
	}
}

// main函数
func main() {
	Listen, err := net.Listen("tcp", "0.0.0.0:8888")
	checkerr(err)
	for {
		conn, err := Listen.Accept()
		checkerr(err)
		go handleConn(conn)
	}
}
