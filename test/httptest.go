package main

import (
	"fmt"
)

func main() {
	tmp := make([]byte, 256)
	tmp = append(tmp[:0])
	if len(tmp) < 0 {
		fmt.Println("xx")
	}
	fmt.Println(tmp)
	// client := &http.Client{Timeout: 5 * time.Second}
	// body := bytes.NewReader([]byte("DataType=PostData&Data=111"))
	// resp, err := client.Post("http://172.16.155.131/test/proxy.php", "application/x-www-form-urlencoded", body)
	// if err != nil {
	// 	panic(err)
	// }
	// data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(data))
	// resp.Body.Close()
}
