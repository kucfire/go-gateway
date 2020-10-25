package main

import (
	"log"
	"net"
)

func main() {
	// listen server
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	if err != nil {
		panic(err)
	}

	// loop to read message
	for {
		// var data [1024]byte
		data := make([]byte, 1024)
		n, addr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Printf("read failed from udp; err : %v \n", err)
			break
		}

		go func() {
			// todo sth
			// reply data
			log.Printf("addr: %v ; data: %v ; count: %v \n", addr, string(data), n)

			_, err := listener.WriteToUDP([]byte("received success!"), addr)
			if err != nil {
				log.Printf("write failed, err : %v \n")
			}
		}()
	}
}
