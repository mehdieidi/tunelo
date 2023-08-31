package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	serverIP := flag.String("i", "127.0.0.1", "Server IP address.")
	serverPort := flag.String("p", "23230", "Server port.")
	flag.Parse()

	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("[error] opening logs file. err: %v\n", err)
	}
	defer logFile.Close()

	if err := godotenv.Load(); err != nil {
		errStr := "[error] loading env: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	serverAddr := *serverIP + ":" + *serverPort

	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		errStr := "[error] creating tcp listener: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("[+] Listening on", serverAddr)

	buf := make([]byte, 4096)

	for {
		conn, err := listener.Accept()
		if err != nil {
			errStr := "[error] connecting: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		for {
			n, err := conn.Read(buf)
			if err != nil {
				errStr := "[error] reading from the connection: " + err.Error()
				fmt.Println(errStr)
				logFile.WriteString(errStr + "\n")
				break
			}

			fmt.Println("received data")

			go handle(buf[:n], secretKey, logFile)
		}

		conn.Close()
	}
}
