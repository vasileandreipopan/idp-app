// means this is a program I can run
package main

// import the packages we need
import (

	// for printing to the console
	"log"

	// for creating a web server
	"net/http"
)

// main is the entry point of the program. when the program is run, this will be the function that is called first
func main() {

	log.Println("Starting server...")

	// create a new HTTP request multiplexer (router), which will handle
	// incoming HTTP requests and route them to the appropriate handler functions
	mux := http.NewServeMux()

	// define a handler function for the /health endpoint. when a GET request is made to /health,
	// this function will be called
	// w is the response writer, which we can use to send a response back to the client
	// r is the incoming HTTP request, which contains information about the request
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		// Write works with raw bytes, so using []byte() to convert the string to bytes
		w.Write([]byte("ok"))
	})

	log.Println("listening on :8888")

	// start the HTTP server on port 8080, using the mux as the handler for incoming requests
	log.Fatal(http.ListenAndServe(":8888", mux))
}
