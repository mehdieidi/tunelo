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

	wsHandler := wsHandler{
		secretKey: secretKey,
		logFile:   logFile,
	}

	http.HandleFunc("/ws", wsHandler.handleWebSocket)

	fmt.Println("[+] Listening on", serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		errStr := "[error] ws listener: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
}
