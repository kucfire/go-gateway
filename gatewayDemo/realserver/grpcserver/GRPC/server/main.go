package main

import (
	"context"
	"flag"
	"fmt"
	pb "gatewayDemo/realserver/grpcserver/GRPC/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const streamingCount = 0

type server struct{}

// 服务端流式处理 ： 服务端读取流数据并将响应发送出去
func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Println("-- ServerStreamingEcho ---")
	fmt.Printf("request received : :v\n", in)

	// Read requests and send responses， 输入一个in， 输出streamingCount个响应
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}

	return nil
}

// 客户端流式处理 ： 客户端持续读取服务端发送过来response，并将结果打印出来
func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Println("--- ClientStreamingEcho ---")

	// Read requests and send response
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("echo last received message")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}

		message = in.Message
		fmt.Printf("request received : %v. building echo\n", in)
		if err != nil {
			return err
		}
	}
	return nil // 防止编译不通过
}

// 双向流式处理 ：
func (s *server) BidirectuibakStreamingEcho(stream pb.Echo_BidirectuibakStreamingEchoServer) error {
	fmt.Println("--- BidirectionalStreamingEcho ---")

	// Read requests and send response
	for {
		in, err := stream.Recv()
		// read last message
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf("client request received %v, sending echo\n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
	return nil
}

// 输入in， 发送输出相同结果
func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("--- UnaryEcho ---\n")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}

	fmt.Println("md : ", md)
	fmt.Printf("request received : %v, sending echo\n", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

// func (s *server) mustEmbedUnimplementedEchoServercompiler() {

// }

var port = flag.Int("port", 50055, "the port to serve on")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
