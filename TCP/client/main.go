package main

import (
	"go-gateway/TCP/unpack"
	"net"
)

func main() {
	conn, errDial := net.Dial("tcp", "localhost:9090")
	if errDial != nil {
		panic(errDial)
	}
	defer conn.Close()
	for i := 0; i < 100; i++ {
		unpack.Encode(conn, "hello world 0!")
	}
}
