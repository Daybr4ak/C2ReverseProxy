package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"sync"

	"github.com/urfave/cli"
	"gopkg.in/ini.v1"
)

var (
	configini string
	wg        sync.WaitGroup
	STATUS    = "reader"
	dataBuf   = make([]byte, 0, 4096)
)

const (
	version = "0.1 dev"
)

// 抛出异常
func checkerr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// 解析config.ini配置文件
func configRun(cli *cli.Context) error {
	cfg, err := ini.Load(cli.String("c"))
	if err != nil {
		fmt.Println("Fail to read file", err)
		os.Exit(1)
	}

	lhostPort := cfg.Section("SERVER").Key("LHOSTPORT").String()
	rhostPort := cfg.Section("SERVER").Key("RHOSTPORT").String()
	BindServer(lhostPort, rhostPort)

	return nil
}

// 端口监听
func BindServer(lhostPort string, rhostPort string) {
	fmt.Println("Server starts listening to port: ", lhostPort)
	lListen, err := net.Listen("tcp", lhostPort)
	checkerr(err)
	// fmt.Println(reflect.TypeOf(lListen))

	fmt.Println("remote listening to port: ", rhostPort)
	rListen, err := net.Listen("tcp", rhostPort)
	checkerr(err)
	// fmt.Println(reflect.TypeOf(rListen))

	wg.Add(2)
	go lconnAccept(lListen)
	go rconnAccept(rListen)
	wg.Wait()

}

// 循环接收127.0.0.1数据
func lconnAccept(lListen net.Listener) {
	defer wg.Done()
	for {
		lconn, err := lListen.Accept()
		checkerr(err)
		go handleConn(lconn, "lconn")
	}
	// }
}

// 循环接收0.0.0.0数据
func rconnAccept(rListen net.Listener) {
	defer wg.Done()
	for {
		rconn, err := rListen.Accept()
		checkerr(err)
		go handleConn(rconn, "rconn")
	}
}

// 获取数据
func handleConn(rconn net.Conn, flag string) {
	if flag == "lconn" {
		for {
			if STATUS == "reader" {
				tmp := make([]byte, 4096)
				data, err := rconn.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					break
				}
				dataBuf = append(dataBuf, tmp[:data]...)
				fmt.Println(flag, tmp[:data], STATUS)
				STATUS = "writer"
			}
		}
	} else if flag == "rconn" {
		for {
			if STATUS == "writer" {
				tmp := make([]byte, 4096)
				data, err := rconn.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					break
				}
				fmt.Println(flag, tmp[:data], STATUS)
				_, err = rconn.Write(dataBuf)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error:", err)
					}
					break
				}
				dataBuf = append(dataBuf[:0])
				// SendData(dataBuf, rconn)
				STATUS = "reader"
			}
		}
	}

}

// 设置
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// main函数
func main() {
	app := cli.NewApp()
	app.Name = "DReverseServer"
	app.Usage = "A ReverseProxy tools"
	app.Version = "0.0.1 dev"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "config.ini",
			Usage:       "config file of DReverseServer",
			Destination: &configini,
		},
	}
	app.Action = configRun
	app.Run(os.Args)
}
