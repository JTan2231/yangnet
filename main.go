package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"time"
)

const (
	heartbeatInterval = 2 * time.Second
)

func main() {
	clientPort := flag.String("clientPort", "3000", "port of this node")
	serverPort := flag.String("serverPort", "", "port to connect to")
	remoteAddress := "127.0.0.1"

	flag.Parse()

	go func() {
		startServer(*clientPort)
	}()

	if *serverPort != "" {
		go func() {
			startClient(remoteAddress + ":" + *serverPort)
		}()
	}

	select {}
}

// Example TCP Client sending heartbeat
func startClient(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to " + serverAddr)
	}
	defer conn.Close()

	ticker := time.NewTicker(heartbeatInterval)
	for {
		select {
		case <-ticker.C:
			_, err := conn.Write([]byte("heartbeat\n"))
			if err != nil {
				log.Fatal("Connection lost:", err)
			}
		}
	}
}

// Example TCP Server handling heartbeat
func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Listening on address 127.0.0.1:" + port)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error:", err)
			}
			return
		}
		log.Println("Received:", message)
	}
}
