package main

import (
	"fmt"
	"go-gateway/basic/unpack"
	"log"
	"net"
)

func main() {
	// simple tcp server
	// listen ip:port
	listener, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		panic(err)
	}

	// accept client request
	// create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept fail,err : %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()

	for {
		bt, err := unpack.Decode(conn)
		if err != nil {
			log.Printf("read from connect fail, err : %v\n", err)
			break
		}
		str := string(bt)
		fmt.Printf("receive from client, data : %v\n", str)
	}
}
