package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	dataBuf = make([]byte, 0, 1046616)
	clients = make(map[string]net.Conn)
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
		tmp := make([]byte, 1046616)
		// dataErr := conn.SetDeadline(time.Now().Add(2 * time.Second))
		// if dataErr != nil {
		// 	fmt.Println("End of getting data")
		// 	continue
		// }
		data, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				// fmt.Println("Read error:", err)
			}
			break
		}
		if string(tmp[:6]) == "TO:GET" {
			if len(dataBuf) > 0 {
				fmt.Println("databuf:")
				fmt.Println("--------------------------------------------------------------")
				fmt.Println(len(dataBuf))
				fmt.Println(dataBuf)
				fmt.Println("--------------------------------------------------------------")
			}
			// _, err = conn.Write(dataBuf)
			_, err = conn.Write([]byte(base64.URLEncoding.EncodeToString(dataBuf)))
			if err != nil {
				if err != io.EOF {
					// fmt.Println("Get Write error:", err)
				}
				break
			}
			dataBuf = append(dataBuf[:0])
			// }
		} else if string(tmp[:10]) == "TO:CONNECT" {
			fmt.Println(string(tmp[:10]))
			_, err = conn.Write([]byte("CONNECT SUCCESS\n"))
			if err != nil {
				if err != io.EOF {
					// fmt.Println("CONNECT Read error:", err)
				}
				break
			}
		} else if string(tmp[:7]) == "TO:SEND" {
			// fmt.Println("Get tmp:", data)
			b64Data, _ := base64.URLEncoding.DecodeString(string(tmp[7:data]))
			// if len(tmp[7:data]) > 0 {
			// 	fmt.Println("Get Data:")
			// 	fmt.Println("--------------------------------------------------------------")
			// 	fmt.Println(string(b64Data))
			// 	fmt.Println("--------------------------------------------------------------")
			// }
			// fmt.Println("b64Data: ")
			// fmt.Println(len(b64Data))
			// fmt.Println(string(b64Data))
			// break
			go BoradCast(clients, b64Data)
			// if len(tmp[7:data]) == 109 {
			// }
		} else {
			dataBuf = append(dataBuf, tmp[:data]...)
		}
		// fmt.Println(len(dataBuf))
		// fmt.Println("tmp:", string(tmp[:data]))
		// fmt.Println("databuf:", string(dataBuf))
	}
}

// 广播消息
func BoradCast(clients map[string]net.Conn, message []byte) {
	for index, client := range clients {
		_, err := client.Write(message)
		if err != nil {
			// fmt.Printf("error: %s\n", err)
			delete(clients, index)
		}
	}
}

// main函数
func main() {
	Listen, err := net.Listen("tcp", "0.0.0.0:8888")
	defer Listen.Close()
	checkerr(err)
	for {
		conn, err := Listen.Accept()
		checkerr(err)
		clients[conn.RemoteAddr().String()] = conn
		go handleConn(conn)
	}
}
