package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	wireguardPort := flag.String(
		"wp",
		"51820",
		"WireGuard port.",
	)
	serverIP := flag.String(
		"i",
		"127.0.0.1",
		"Server IP address.",
	)
	serverPort := flag.String(
		"p",
		"23230",
		"Server port.",
	)
	flag.Parse()

	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("[error] opening logs file. err: %v\n", err)
	}
	defer logFile.Close()

	if err := godotenv.Load(); err != nil {
		errStr := fmt.Sprintf("[error] loading env: %v\n", err.Error())
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if string(secretKey) == "" {
		errStr := "[error] secret key cannot be empty\n"
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}

	wgAddr := "127.0.0.1:" + *wireguardPort

	handler := handler{
		wgAddr:    wgAddr,
		secretKey: secretKey,
		logFile:   logFile,
	}

	http.HandleFunc("/ws", handler.wsHandler)

	serverAddr := *serverIP + ":" + *serverPort
	fmt.Println("[+] HTTP server listening to", serverAddr)

	if err = http.ListenAndServe(serverAddr, nil); err != nil {
		errStr := fmt.Sprintf("[error] http listener: %v\n", err.Error())
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}
}
