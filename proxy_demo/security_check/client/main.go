package main

import (
	"fmt"
	"github.com/donscoco/gateway/base_server/jwt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	// 创建连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //探活时间
		}).DialContext,
		MaxIdleConns:          100,              //最大空闲连接
		IdleConnTimeout:       90 * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
	}

	// 创建客户端
	client := &http.Client{
		Timeout:   time.Second * 30, //请求超时时间
		Transport: transport,
	}

	// 请求数据
	req, err := http.NewRequest("GET", "http://127.0.0.1:8001", nil)
	if err != nil {
		panic(err)
	}

	token, _ := jwt.Encode("ironhead")

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bds, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bds))
}
