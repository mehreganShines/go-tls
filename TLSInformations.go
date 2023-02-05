//incomplete but works. more features will add soon
// A simple Go program that can print information about the TLS handshakes with itself and listening on port 8443.
//to generate certificates run "go-tls/self-cert.go" and save the results in cert.pem and key.pem files.
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading cert and key:", err)
		os.Exit(1)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on port 8443...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	tlsConn := conn.(*tls.Conn)

	err := tlsConn.Handshake()
	if err != nil {
		fmt.Println("Error performing handshake:", err)
		return
	}

	state := tlsConn.ConnectionState()
	fmt.Println("TLS Version:", state.Version)
	fmt.Println("Cipher Suite:", state.CipherSuite)
	fmt.Println("Peer Certificates:")
	for _, cert := range state.PeerCertificates {
		fmt.Println("- Subject:", cert.Subject)
		fmt.Println("  Issuer:", cert.Issuer)
	}
}
