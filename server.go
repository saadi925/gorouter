package flow

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is a wrapper around http.Server.
type Server struct {
	*http.Server
}

// ServerConfig represents configuration options for the flow server.
type ServerConfig struct {
	Addr         string        // Server address (e.g., ":8080")
	ReadTimeout  time.Duration // Read timeout for incoming requests
	WriteTimeout time.Duration // Write timeout for outgoing responses
	IdleTimeout  time.Duration // Idle timeout for keep-alive connections
	TLSConfig    TLSConfig     // TLS/SSL configuration
}

// TLSConfig represents TLS/SSL configuration options.
type TLSConfig struct {
	CertFile string // Path to the SSL certificate file
	KeyFile  string // Path to the SSL private key file
}

// NewServer creates a new instance of the flow server with the given configuration.
func NewServer(handler http.Handler, config ServerConfig) *Server {
	server := &http.Server{
		Addr:         config.Addr,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	if config.TLSConfig.CertFile != "" && config.TLSConfig.KeyFile != "" {
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		var err error
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(config.TLSConfig.CertFile, config.TLSConfig.KeyFile)
		if err != nil {
			log.Fatalf("Error loading TLS certificate: %v", err)
		}
		server.TLSConfig = tlsConfig
	}

	return &Server{
		Server: server,
	}
}

// Start starts the flow server.
func (s *Server) Start() error {
	log.Printf("Server starting on address %s...", s.Addr)
	if s.TLSConfig != nil {
		return s.ListenAndServeTLS("", "")
	}
	return s.ListenAndServe()
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	DSN string
}

func LoadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Stop stops the flow server gracefully.
func (s *Server) Stop() error {
	log.Println("Server shutting down gracefully...")
	return s.Shutdown(nil)
}

// GracefulShutdown handles graceful shutdown of the server.
func GracefulShutdown(server *Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh // Wait for termination signal

	log.Println("Shutting down server...")

	if err := server.Stop(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped")
}
