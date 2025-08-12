package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	jplaw "go.ngs.io/jplaw-api-v2"
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
	http.HandleFunc("/epubs/", epubsHandler)
	http.HandleFunc("/graphql", graphqlHandler)
	http.HandleFunc("/graphiql", graphiqlHandler)

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

func epubsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path: /epubs/{id}
	path := strings.TrimPrefix(r.URL.Path, "/epubs/")
	if path == "" || path == "/" {
		http.Error(w, "Law ID is required", http.StatusBadRequest)
		return
	}

	// Remove any trailing slashes
	lawID := strings.TrimSuffix(path, "/")

	log.Printf("Fetching law data for ID: %s", lawID)

	// Create API client
	client := jplaw.NewClient()

	// Set up parameters to get XML format
	xmlFormat := jplaw.ResponseFormatXml
	params := &jplaw.GetLawDataParams{
		LawFullTextFormat: &xmlFormat,
	}

	// Get law data with XML format
	lawData, err := client.GetLawData(lawID, params)
	if err != nil {
		log.Printf("Error fetching law data for ID %s: %v", lawID, err)
		if strings.Contains(err.Error(), "404") {
			http.Error(w, "Law not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching law data: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Extract XML content from response
	var xmlContent []byte
	if lawData.LawFullText != nil {
		// The LawFullText should contain the XML as a string when LawFullTextFormat is XML
		if xmlStr, ok := (*lawData.LawFullText).(string); ok {
			// The XML is Base64 encoded, decode it
			decodedXML, err := base64.StdEncoding.DecodeString(xmlStr)
			if err != nil {
				log.Printf("Error decoding Base64 for law ID %s: %v", lawID, err)
				http.Error(w, "Error decoding XML content", http.StatusInternalServerError)
				return
			}
			// Remove <TmpRootTag> wrapper if present
			xmlStr := string(decodedXML)
			if strings.HasPrefix(xmlStr, "<TmpRootTag>") {
				xmlStr = strings.TrimPrefix(xmlStr, "<TmpRootTag>")
				xmlStr = strings.TrimSuffix(xmlStr, "</TmpRootTag>")
			}
			xmlContent = []byte(xmlStr)
			log.Printf("Decoded XML content length for law ID %s: %d bytes", lawID, len(xmlContent))
		} else {
			log.Printf("Unexpected type for LawFullText: %T", *lawData.LawFullText)
			http.Error(w, "Invalid XML format in response", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("No LawFullText in response for law ID %s", lawID)
		http.Error(w, "No law content in response", http.StatusInternalServerError)
		return
	}

	// Convert XML to EPUB
	xmlReader := bytes.NewReader(xmlContent)
	book, err := jplaw2epub.CreateEPUBFromXMLFile(xmlReader)
	if err != nil {
		log.Printf("Error creating EPUB for law ID %s: %v", lawID, err)
		http.Error(w, fmt.Sprintf("Error creating EPUB: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate EPUB to buffer
	var buf bytes.Buffer
	if _, err := book.WriteTo(&buf); err != nil {
		log.Printf("Error writing EPUB to buffer for law ID %s: %v", lawID, err)
		http.Error(w, "Error generating EPUB", http.StatusInternalServerError)
		return
	}

	// Set response headers
	filename := fmt.Sprintf("%s.epub", lawID)
	w.Header().Set("Content-Type", "application/epub+zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	// Write response
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Error writing response for law ID %s: %v", lawID, err)
		return
	}

	log.Printf("Successfully converted law ID %s to EPUB (%d bytes)", lawID, buf.Len())
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var query string
	var variables map[string]interface{}

	if r.Method == http.MethodGet {
		query = r.URL.Query().Get("query")
		variablesStr := r.URL.Query().Get("variables")
		if variablesStr != "" {
			json.Unmarshal([]byte(variablesStr), &variables)
		}
	} else if r.Method == http.MethodPost {
		var body struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		query = body.Query
		variables = body.Variables
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         SchemaFixed,
		RequestString:  query,
		VariableValues: variables,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func graphiqlHandler(w http.ResponseWriter, r *http.Request) {
	h := handler.New(&handler.Config{
		Schema:   &SchemaFixed,
		Pretty:   true,
		GraphiQL: true,
	})

	h.ServeHTTP(w, r)
}
