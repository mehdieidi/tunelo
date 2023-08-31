package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/chacha20poly1305"
)

func main() {
	wireguardPort := flag.String(
		"wp",
		"23232",
		"The port that the Wireguard is listening to.",
	)
	listeningPort := flag.String(
		"p",
		"23231",
		"The port that the app is listening to.",
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

	fmt.Println("Wireguard port:", *wireguardPort)
	fmt.Println("App listening on port:", *listeningPort)
	fmt.Println("Remote server IP/Port :", *remoteServerIP, *remoteServerPort)

	go gracefulShutdown()

	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("[error] opening logs file. err: %v\n", err)
	}
	defer logFile.Close()

	serverAddr := *remoteServerIP + ":" + *remoteServerPort

	conn, err := net.ListenPacket("udp", "127.0.0.1:"+(*listeningPort))
	if err != nil {
		errStr := "[error] listening on wireguard interface: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[+] Running UDP listener on port", *listeningPort)

	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		errStr := "[error] connecting to the remote server: " + err.Error()
		fmt.Println(errStr)
		logFile.WriteString(errStr + "\n")
		os.Exit(1)
	}
	defer serverConn.Close()

	secretKey := []byte("123abc@#$456asdfg#$%89756*&^fegv")

	buf := make([]byte, 2048)

	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			errStr := "[error] reading data from the wireguard interface: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		fmt.Println("received from udp conn:", string(buf))

		nonce := make([]byte, chacha20poly1305.NonceSize)
		if _, err := rand.Read(nonce); err != nil {
			errStr := "[error] generating nonce: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		aead, err := chacha20poly1305.New(secretKey)
		if err != nil {
			errStr := "[error] creating aead: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}
		encryptedData := aead.Seal(nil, nonce, buf[:n], nil)

		encryptedDataWithNonce := append(nonce, encryptedData...)

		_, err = serverConn.Write(encryptedDataWithNonce)
		if err != nil {
			errStr := "[error] writing to the server tcp connection: " + err.Error()
			fmt.Println(errStr)
			logFile.WriteString(errStr + "\n")
			continue
		}

		logFile.WriteString("[info] data just flew away!" + "\n")
	}
}

func gracefulShutdown() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		fmt.Println("\n[-] shutdown signal received")

		defer func() {
			stop()
		}()

		os.Exit(0)
	}()
}
