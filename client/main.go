package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"

	"donatello/pkg/graceful"
)

func main() {
	wireGuardPort := flag.String(
		"wp",
		"51820",
		"WireGuard port.",
	)
	listeningPort := flag.String(
		"lp",
		"23231",
		"Local port that the app is listening to.",
	)
	remoteServerIP := flag.String(
		"ri",
		"127.0.0.1",
		"Remote server IP.",
	)
	remoteServerPort := flag.String(
		"rp",
		"23230",
		"Remote server port.",
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

	udpServerAddr := "127.0.0.1:" + *listeningPort
	udpListener, err := net.ListenPacket("udp", udpServerAddr)
	if err != nil {
		errStr := fmt.Sprintf("[error] creating udp listener: %v\n", err.Error())
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}
	defer udpListener.Close()

	fmt.Println("[+] UDP server listening to", udpServerAddr)

	wsServerAddr := fmt.Sprintf("ws://%s:%s/ws", *remoteServerIP, *remoteServerPort)
	wsConn, err := connectServerWS(wsServerAddr)
	if err != nil {
		errStr := fmt.Sprintf("[error] ws dial: %v\n", err.Error())
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}
	defer wsConn.Close()

	fmt.Println("[+] WS connected:", wsServerAddr)

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	if string(secretKey) == "" {
		errStr := "[error] secret key cannot be empty."
		fmt.Println(errStr)
		logFile.WriteString(errStr)
		os.Exit(1)
	}

	wgAddr := "127.0.0.1" + ":" + *wireGuardPort

	handler := handler{
		wsServerAddr: wsServerAddr,
		wgAddr:       wgAddr,
		wsConn:       wsConn,
		udpListener:  udpListener,
		secretKey:    secretKey,
		logFile:      logFile,
	}

	go handler.wsReadHandler()
	go handler.udpReadHandler()

	graceful.ShutdownHandler()
}
