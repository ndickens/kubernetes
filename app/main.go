package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// var logFilePath = "/Users/ndickens/Projects/kubernetes-project/app/ip_addresses.log" // For local testing
var logFilePath = "/logs/ip_addresses.log" // Path to the log file inside the Persistent Volume

func helloWorldHandler(w http.ResponseWriter, req *http.Request) {
	// Generate the log message
	ip := req.RemoteAddr
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf("[INFO] [%s] IP address: %s\n", timestamp, ip)
	fmt.Printf(logMessage)

	// Open the log file in append mode and create if it doesn't exist
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	// Write the IP address to the log file
	if _, err := f.WriteString(logMessage); err != nil {
		log.Println("Error writing to log file:", err)
		return
	}

	// Return "Hello World"
	fmt.Fprintf(w, "Hello, world!")
}

func logHandler(w http.ResponseWriter, req *http.Request) {
	// Read the log file
	content, error := os.ReadFile(logFilePath)
	if error != nil {
		fmt.Fprintf(w, "File does not exist")
		return
	}

	// Return file contents
	fmt.Fprintf(w, string(content))
}

func main() {
	// Define the route and handler for "/hello"
	http.HandleFunc("/hello", helloWorldHandler)

	// Created for testing purposes to quickly verify that log contents persist
	http.HandleFunc("/logs", logHandler)

	// Start the server on port 8080
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
