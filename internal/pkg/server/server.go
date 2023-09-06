// server.go
package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

// Server represents your gRPC server.
type Server struct {
    grpcServer *grpc.Server
    wg         sync.WaitGroup
    stopped    bool // Custom flag to track server status
}

// NewServer creates a new instance of the Server.
func NewServer() *Server {
    return &Server{
        grpcServer: grpc.NewServer(),
    }
}

// Start starts the gRPC server.
func (s *Server) Start(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    fmt.Printf("Server listening on port %s...\n", port)

    // Register your gRPC service implementation here, e.g., auth.RegisterUserServiceServer(s.grpcServer, &yourService{})

    // Start the server in a goroutine
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        if err := s.grpcServer.Serve(listener); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()
}

// Stop stops the gRPC server gracefully.
func (s *Server) Stop() {
    // Gracefully stop the gRPC server
    s.grpcServer.GracefulStop()

    // Set the custom flag to indicate that the server has stopped
    s.stopped = true
}

// Wait waits for the server to finish gracefully.
func (s *Server) Wait() {
    // Create a signal channel to capture termination signals (e.g., SIGINT, SIGTERM)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // Wait for termination signals or server shutdown
    select {
    case <-sigCh:
        // Handle termination signals if needed
    case <-s.waitServerStopped():
        // Server has stopped gracefully
    }

    // Wait for any remaining goroutines to finish
    s.wg.Wait()
}

// waitServerStopped waits until the server is marked as stopped.
func (s *Server) waitServerStopped() <-chan struct{} {
    ch := make(chan struct{})
    go func() {
        for !s.stopped {
            // Wait until the server is stopped
        }
        close(ch)
    }()
    return ch
}
