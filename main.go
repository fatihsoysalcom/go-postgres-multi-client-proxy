package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	proxyPort   = flag.String("port", "6000", "Port for the proxy to listen on")
	backendAddr = flag.String("backend", "localhost:5432", "Address of the backend Postgres database (e.g., localhost:5432)")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*proxyPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *proxyPort, err)
	}
	defer listener.Close()

	log.Printf("Go Postgres MCP Proxy listening on :%s, forwarding to %s", *proxyPort, *backendAddr)
	log.Println("Ensure a Postgres database is running at the backend address.")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}
		log.Printf("Accepted client connection from %s", clientConn.RemoteAddr())
		// Handle each client connection in a new goroutine to support multiple concurrent clients.
		go handleClient(clientConn, *backendAddr)
	}
}

// handleClient establishes a connection to the backend and proxies data bidirectionally.
func handleClient(clientConn net.Conn, backendAddr string) {
	// Ensure client connection is closed when this function exits.
	// This is crucial for preventing connection leaks from the proxy to the client side.
	defer func() {
		log.Printf("Closing client connection from %s", clientConn.RemoteAddr())
		clientConn.Close()
	}()

	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		log.Printf("Failed to connect to backend %s for client %s: %v", backendAddr, clientConn.RemoteAddr(), err)
		return // Client connection will be closed by defer
	}
	// Ensure backend connection is closed when this function exits.
	// This prevents connection leaks from the proxy to the backend database.
	defer func() {
		log.Printf("Closing backend connection to %s for client %s", backendAddr, clientConn.RemoteAddr())
		backendConn.Close()
	}()

	log.Printf("Proxying connection %s <-> %s <-> %s", clientConn.RemoteAddr(), clientConn.LocalAddr(), backendConn.RemoteAddr())

	// Use a channel to wait for either copy operation to finish.
	// This ensures that when one side closes, the other side's copy goroutine
	// is also effectively shut down by the deferred connection closes.
	done := make(chan struct{}, 2) // Buffered channel to avoid blocking if both goroutines finish quickly

	// Copy data from client to backend
	go func() {
		_, err := io.Copy(backendConn, clientConn)
		if err != nil && err != io.EOF { // io.EOF is expected on graceful close
			log.Printf("Error copying client to backend for %s: %v", clientConn.RemoteAddr(), err)
		}
		done <- struct{}{} // Signal that this copy operation is done
	}()

	// Copy data from backend to client
	go func() {
		_, err := io.Copy(clientConn, backendConn)
		if err != nil && err != io.EOF { // io.EOF is expected on graceful close
			log.Printf("Error copying backend to client for %s: %v", clientConn.RemoteAddr(), err)
		}
		done <- struct{}{} // Signal that this copy operation is done
	}()

	// Wait for one of the copy operations to complete.
	// When one side closes its connection, its io.Copy will return,
	// sending a signal to 'done'. This allows the handleClient function
	// to exit, triggering the deferred closes for both client and backend connections.
	<-done
	log.Printf("Proxy connection for %s finished.", clientConn.RemoteAddr())
	// The other io.Copy goroutine will eventually fail with a "use of closed network connection"
	// or similar error when its peer connection is closed by the defer, or it will naturally
	// return if the peer also closed gracefully.
}
