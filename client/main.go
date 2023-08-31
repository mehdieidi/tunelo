package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	listeningPort := flag.String(
		"p",
		"23231",
		"The local port that the app is listening to.",
	)
	remoteServerIP := flag.String(
		"ri",
		"127.0.0.1",
		"The remote server IP.",
	)
	remoteServerPort := flag.String(
		"rp",
		"23230",
		"The remote server port.",
	)
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

	go handleGracefulShutdown()

	wsServerAddr := fmt.Sprintf("ws://%s:%s/ws", *remoteServerIP, *remoteServerPort)
	udpServerAddr := "127.0.0.1:" + *listeningPort

	udpListener, err := net.ListenPacket("udp", udpServerAddr)
	if err != nil {
		errStr := "[error] creating udp listener: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
	defer udpListener.Close()

	fmt.Println("[+] UDP listener on", udpServerAddr)

	wsConn, _, err := websocket.DefaultDialer.Dial(wsServerAddr, nil)
	if err != nil {
		errStr := "[error] ws dial: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
	defer wsConn.Close()

	buf := make([]byte, 4096)

	for {
		n, _, err := udpListener.ReadFrom(buf)
		if err != nil {
			errStr := "[error] reading data from the udp conn: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		fmt.Println("received from udp conn:", string(buf))

		encryptedData, err := encrypt(buf[:n], secretKey)
		if err != nil {
			fmt.Println(err)
			logFile.WriteString(err.Error() + "\n")
			continue
		}

		err = wsConn.WriteMessage(websocket.BinaryMessage, encryptedData)
		if err != nil {
			errStr := "[error] ws write: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		logFile.WriteString("[info] data just flew away!" + "\n")
	}
}
