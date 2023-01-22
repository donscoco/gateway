package main

import (
	"flag"
	"fmt"
	"github.com/donscoco/gateway/base_server/zookeeper"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", "localhost:8002", "http service address")
var addr2 = flag.String("addr2", "localhost:8003", "http service address")

func main() {

	flag.Parse()
	log.SetFlags(0)

	if len(*addr) > 0 {
		server1 := MyServer{Addr: *addr}
		server1.Run()
	}
	if len(*addr2) > 0 {
		server2 := MyServer{Addr: *addr2}
		server2.Run()
	}

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type MyServer struct {
	Addr string
}

func (r *MyServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_http_string/test_http_string/aaa", r.TimeoutHandler)

	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		//注册zk节点
		zkManager := zookeeper.NewZkManager([]string{"192.168.2.132:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Printf(" connect zk error: %s ", err)
		}
		defer zkManager.Close()
		err = zkManager.RegistServerPath("/real_server", r.Addr)
		if err != nil {
			fmt.Printf(" regist node error: %s ", err)
		}
		zlist, err := zkManager.GetServerListByPath("/real_server")
		fmt.Println(zlist)
		log.Fatal(server.ListenAndServe())
	}()

}

func (r *MyServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	//127.0.0.1:8008/abc?sdsdsa=11
	//r.Addr=127.0.0.1:8008
	//req.URL.Path=/abc
	//fmt.Println(req.Host)
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("headers =%v\n", req.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, header)
	io.WriteString(w, r.Addr+"\n")

}

func (r *MyServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *MyServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}
