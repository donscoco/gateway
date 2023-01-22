// Binary server is an example server.
package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/donscoco/gateway/server_demo/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 30801, "the port to serve on")

const (
	streamingCount = 10
)

// 去实现 proto 文件中定义的四种接口
type server struct{}

// 1. 简单RPC（Simple RPC）：即客户端发送一个请求给服务端，从服务端获取一个应答，就像一次普通的函数调用。
func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("--- UnaryEcho ---\n")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Println("md", md)
	fmt.Printf("request received: %v, sending echo\n", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

// 2. 服务端流式RPC（Server-side streaming RPC）：一个请求对象，服务端可以传回多个结果对象。即客户端发送一个请求给服务端，可获取一个数据流用来读取一系列消息。客户端从返回的数据流里一直读取直到没有更多消息为止。
func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("--- ServerStreamingEcho ---\n")
	fmt.Printf("request received: %v\n", in)
	// Read requests and send responses.
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

// 3. 客户端流式RPC（Client-side streaming RPC）：客户端传入多个请求对象，服务端返回一个响应结果。即客户端用提供的一个数据流写入并发送一系列消息给服务端。一旦客户端完成消息写入，就等待服务端读取这些消息并返回应答。
func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("--- ClientStreamingEcho ---\n")
	// Read requests and send responses.
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

// 4. 双向流式RPC（Bidirectional streaming RPC）：结合客户端流式rpc和服务端流式rpc，可以传入多个对象，返回多个响应对象。即两边都可以分别通过一个读写数据流来发送一系列消息。这两个数据流操作是相互独立的，所以客户端和服务端能按其希望的任意顺序读写，例如：服务端可以在写应答前等待所有的客户端消息，或者它可以先读一个消息再写一个消息，或者是读写相结合的其他方式。每个数据流里消息的顺序会被保持。
func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Printf("--- BidirectionalStreamingEcho ---\n")
	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
