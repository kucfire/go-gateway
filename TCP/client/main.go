package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// 连接服务器
	conn, errDial := net.Dial("tcp", "localhost:9090")
	defer conn.Close()
	if errDial != nil {
		log.Printf("connect failed, err : %v \n", errDial)
		return
		// panic(errDial)
	}

	// for i := 0; i < 100; i++ {
	// 	unpack.Encode(conn, "hello world 0!")
	// }
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 一直读取，直到读到\n
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Printf("read from consolo failed, err : %v\n", err)
			break
		}

		trimmedInput := strings.TrimSpace(input)

		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			log.Printf("write failed, err : %v\n", err)
			break
		}
		// time.Sleep(10 * time.Second)

		if trimmedInput == "Q" {
			log.Println("cilent exit!")
			break
		}
	}
}
