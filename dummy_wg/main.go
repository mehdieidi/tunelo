package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:23231")
	if err != nil {
		panic(err)
	}

	var i int

	for {
		time.Sleep(1 * time.Second)

		_, err := conn.Write([]byte("dummy data" + strconv.Itoa(i)))
		if err != nil {
			fmt.Println("error sending data", err)
			continue
		}

		fmt.Println("sent dummy data", i)

		i++
	}
}
