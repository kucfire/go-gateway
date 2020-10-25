package main

import (
	"log"
	"net"
)

func main() {
	// create conn
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	if err != nil {
		log.Printf("connect fail, err : %v \n", err)
		panic(err)
	}

	for i := 0; i < 100; i++ {
		// send data
		_, err := conn.Write([]byte("hello,server!"))
		if err != nil {
			log.Printf("send data failed, err : %v \n", err)
			return
		}

		// receive data
		result := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(result)
		if err != nil {
			log.Printf("receive data failed, err : %v \n", err)
			return
		}
		log.Printf("receive from addr : %v, data : %v \n", remoteAddr, string(result[:n]))

	}
}
