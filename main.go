package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// The release of the package. This is supplied by the compiler.
var release string

// The internally-used version of the package. This allows for a fallback if
// the version number is not defined, such as in tests or what not.
var version string

func init() {
	if release != "" {
		version = release
	} else {
		version = "0.0.1-UNRELEASED"
	}
}

// helloMsg holds a string tempalte for our hello world! Message on the default
// route.
const helloMsg = `
<!DOCTYPE HTML>
<html>
<head>
<title>Hello World!</title>
</head>
Hello %s!<br>
<br>
From host %s on version %s<br>
<body>
</body>
</html>
`

// clientAddress fetches the cilent address from the HTTP request. It first
// tries X-Forwarded-For - if that header does not exist, it gets the address
// from RemoteAddr.
func clientAddress(r *http.Request) (client string) {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		client = strings.Split(v, ", ")[0]
	} else {
		var err error
		client, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error getting host from remote address: %s", err)
		}
	}
	return
}

// logRequest is a helper to log all incoming requests.
func logRequest(r *http.Request, code int) {
	client := clientAddress(r)

	log.Printf("%s %s - %d - %s - %s", r.Method, r.URL.Path, code, client, r.UserAgent())
}

// handleVersion prints the version of the package in JSON format.
func handleVersion(w http.ResponseWriter, r *http.Request) {
	logRequest(r, 200)
	w.Header().Add("Content-Type", "application/json")
	if _, err := io.WriteString(w, fmt.Sprintf("{\"version\": \"%s\"}", version)); err != nil {
		log.Printf("Error writing version JSON string to client: %s", err)
	}
}

// handleDefault prints a "hello world!" message.
func handleDefault(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		hostName, err := os.Hostname()
		if err != nil {
			logRequest(r, 500)
			log.Printf("Error getting hostname: %s", err)
			http.Error(w, fmt.Sprintf("Error getting hostname: %s", err), 500)
			return
		}
		logRequest(r, 200)
		client := clientAddress(r)
		if _, err := io.WriteString(w, fmt.Sprintf(helloMsg, client, hostName, version)); err != nil {
			log.Printf("Error writing default response to client: %s", err)
		}
		return
	}
	logRequest(r, 404)
	http.Error(w, fmt.Sprintf("Path %s not found", r.URL.Path), 404)
}

// startServer starts the HTTP server and runs it until the listening socket is
// closed, or if an interrupt or SIGTERM is sent.
func startServer(ln net.Listener) error {
	http.HandleFunc("/version", handleVersion)
	http.HandleFunc("/", handleDefault)
	log.Printf("Press CTRL-C or send SIGTERM to close the server")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Printf("%s received, shutting down.", s.String())
		if err := ln.Close(); err != nil {
			log.Printf("Error closing listener: %s", err)
		}
		os.Exit(0)
	}()
	return http.Serve(ln, nil)
}

func main() {
	log.Println("Server starting...")

	// #nosec
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error listening on %s: %s", ln, err)
	}
	log.Printf("Listening on %s", ln.(*net.TCPListener).Addr().String())

	log.Fatal(startServer(ln))
}
