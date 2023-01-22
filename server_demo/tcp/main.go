package main

import (
	"context"
	"fmt"
	"github.com/donscoco/gateway/base_server/zookeeper"
	server_tcp "github.com/donscoco/gateway/server/tcp"
	"github.com/donscoco/gateway/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type MyServer struct {
	Addr string
}

func (s *MyServer) Run() {
	// 1.创建 handler 实现 ServeTCP 方法以实现 TCPHandler 接口
	// 2.创建 TCP 服务器，注册进上面的handler
	// 3.listenAndServe

	tcphandler := &tcpHandler{}
	tcpServer := &server_tcp.TcpServer{
		Addr:    s.Addr,
		Handler: tcphandler,
	}
	tcphandler.Server = tcpServer
	go func() {
		//注册zk节点
		zkManager := zookeeper.NewZkManager([]string{"192.168.2.132:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Printf(" connect zk error: %s ", err)
		}
		defer zkManager.Close()
		err = zkManager.RegistServerPath("/real_server", s.Addr)
		if err != nil {
			fmt.Printf(" regist node error: %s ", err)
		}
		zlist, err := zkManager.GetServerListByPath("/real_server")
		fmt.Println(zlist)
		log.Println("listen on ", s.Addr)
		log.Fatal(tcpServer.ListenAndServe())
	}()

}

type tcpHandler struct {
	Server *server_tcp.TcpServer
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte(t.Server.Addr + ":tcpHandler\n"))
}

func main() {

	// 1.创建 handler 实现 ServeTCP 方法以实现 TCPHandler 接口
	// 2.创建 TCP 服务器，注册进上面的handler
	// 3.listenAndServe

	innerNet := util.GetIpv4_192_168()
	server1 := MyServer{Addr: innerNet + ":7010"}
	server1.Run()
	server2 := MyServer{Addr: innerNet + ":7020"}
	server2.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	//代理测试
	//rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	//rb.Add("127.0.0.1:6001", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//redis服务器测试
	//rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	//rb.Add("127.0.0.1:6379", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//http服务器测试:
	//缺点对请求的管控不足,比如我们用来做baidu代理,因为无法更改请求host,所以很轻易把我们拒绝
	//rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	//rb.Add("127.0.0.1:2003", "40")
	////rb.Add("www.baidu.com:80", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("tcp_proxy start at:" + addr)
	//tcpServ.ListenAndServe()

	//websocket服务器测试:缺点对请求的管控不足
	//rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	//rb.Add("127.0.0.1:2003", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//http2服务器测试:缺点对请求的管控不足
	//rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	//rb.Add("127.0.0.1:3003", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()
}
