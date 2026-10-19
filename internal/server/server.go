// Package server provides a simple HTTP server for serving generated
// documentation locally during development.
package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Server serves static files from a directory.
type Server struct {
	RootDir string
	Port    int
	server  *http.Server
}

// New creates a server that serves files from rootDir on the given port.
func New(rootDir string, port int) *Server {
	return &Server{
		RootDir: rootDir,
		Port:    port,
	}
}

// Start begins serving files.  Blocks until the server is stopped.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}

	return s.server.ListenAndServe()
}

// Stop gracefully stops the server.
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}

// Address returns the address the server will listen on.
func (s *Server) Address() string {
	return fmt.Sprintf(":%d", s.Port)
}

// handleRequest serves files from the root directory, with a fallback
// to index.html for directory requests.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(r.URL.Path)
	if strings.HasPrefix(cleanPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	filePath := filepath.Join(s.RootDir, cleanPath)

	// If the path is a directory, serve index.html
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Set content type based on extension
	setContentType(w, filePath)

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// setContentType sets the Content-Type header based on file extension.
func setContentType(w http.ResponseWriter, path string) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		// Let http.ServeFile detect the content type
	}
}

// IsPortAvailable checks whether the given port is available for binding.
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
