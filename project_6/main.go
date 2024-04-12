package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"reliable-tcp/ftp"
)

func main() {
	var port int
	var dst string
	var isServer bool

	flag.IntVar(&port, "port", 2000, "Port number to run the server")
	flag.StringVar(&dst, "dst", "", "Destination address")
	flag.BoolVar(&isServer, "s", false, "Run as server")

	flag.Parse()

	if isServer {
		// Create a UDP listener
		addr := fmt.Sprintf(":%d", port)
		serverAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.ListenUDP("udp", serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Handle incoming requests
		for {
			handleClient(conn)
		}
	} else {
		// Client mode
		ftp.Run(dst, port)
	}
}

func handleClient(conn *net.UDPConn) {
	// Receive data from client
	buffer := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	// Process the received data (in this case, assume it's a file path)
	data := string(buffer[:n])

	log.Printf("Received request from %s: %s\n", addr.String(), data)

	// Open the requested file
	file, err := os.Open(data)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file contents
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Error getting file information:", err)
		return
	}
	fileSize := fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	_, err = file.Read(fileBuffer)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	// Send the file back to the client
	_, err = conn.WriteToUDP(fileBuffer, addr)
	if err != nil {
		log.Println("Error sending file to client:", err)
		return
	}

	log.Println("File sent successfully to", addr.String())
}
