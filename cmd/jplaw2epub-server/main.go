package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"go.ngs.io/jplaw2epub"
)

func main() {
	portFlag := flag.String("port", "", "Port to listen on (default: find available port)")
	flag.Parse()

	var port string
	if *portFlag != "" {
		port = *portFlag
	} else if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	} else {
		// Find an available port
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("Failed to find available port: %v", err)
		}
		port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
		listener.Close()
	}

	http.HandleFunc("/convert", convertHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/xml" && contentType != "text/xml" {
		http.Error(w, "Content-Type must be application/xml or text/xml", http.StatusBadRequest)
		return
	}

	xmlData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if len(xmlData) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	xmlReader := bytes.NewReader(xmlData)
	book, err := jplaw2epub.CreateEPUBFromXMLFile(xmlReader)
	if err != nil {
		log.Printf("Error creating EPUB: %v", err)
		http.Error(w, fmt.Sprintf("Error creating EPUB: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if _, err := book.WriteTo(&buf); err != nil {
		log.Printf("Error writing EPUB to buffer: %v", err)
		http.Error(w, "Error generating EPUB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/epub+zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"law.epub\"")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	log.Printf("Successfully converted XML to EPUB (%d bytes)", buf.Len())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","service":"jplaw2epub-server"}`)
}