package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// simple tcp server
	// listen ip:port
	listener, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Printf("listen failed, err : %v \n", err)
		return
		// panic(err)
	}

	// accept client request
	// create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept fail,err : %v\n", err)
			continue
		}
		// closeMgs := make(chan bool)
		go process(conn)
		// if <-closeMgs {
		// 	log.Println("server exit!")
		// 	break
		// }
	}
}

func process(conn net.Conn) {
	// defer conn.Close()

	for {
		// bt, err := unpack.Decode(conn)
		buf := make([]byte, 128)
		_, err := conn.Read(buf)
		if err != nil {
			log.Printf("read from connect fail, err : %v\n", err)
			break
		}

		str := string(buf)

		fmt.Printf("receive from client, data : %v\n", str)
	}
}
