// means this is a program I can run
package main

// import the packages we need
import (

	// convert Go data into JSON.
	"encoding/json"

	// for printing to the console
	"log"

	// for creating a web server
	"net/http"

	// for getting environment variables
	"os"

	// for setting timeouts on the server
	"time"
)

// package-level variable, every function in the file can use this
// "dev" is the default value
var version = "dev"

// w is the response writer, which is used to send data back to the client
// status is the HTTP status code to send back to the client
// v is the data to send back to the client, which will be converted to JSON
func writeJSON(w http.ResponseWriter, status int, v interface{}) {

	// set the Content-Type header to application/json, so the client knows we're sending JSON data
	w.Header().Set("Content-Type", "application/json")

	// set the HTTP status code to send back to the client
	w.WriteHeader(status)

	// convert the data to JSON and send it back to the client
	// _ is used to ignore the error returned by json.NewEncoder(w).Encode(v), because we don't need to handle it in this case
	_ = json.NewEncoder(w).Encode(v)
}

// main is the entry point of the program. when the program is run, this will be the function that is called first
func main() {

	log.Println("Starting server...")

	// create a new HTTP request multiplexer (router), which will handle
	// incoming HTTP requests and route them to the appropriate handler functions
	mux := http.NewServeMux()

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"version": version})
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// get the port from the environment variable PORT, or use 8888 if it's not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	// create a new HTTP server with the specified address, handler, and read timeout
	// & because we want to create a pointer to the http.Server struct, so we can modify its fields later if needed
	// ReadTimeout is the Slowloris protection, it will close the connection if the client takes too long to send the request
	srv := &http.Server{Addr: ":" + port, Handler: mux, ReadTimeout: 5 * time.Second}

	log.Printf("idp-app %s listening on :%s", version, port)
	log.Fatal(srv.ListenAndServe())
}
