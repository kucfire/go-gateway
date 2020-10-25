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
	unpack.Encode(conn, "hello world 0!")
}
