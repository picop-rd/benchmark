package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	port = flag.String("port", "80", "listen port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", *port),
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("serve http")
	log.Fatal(server.Serve(ln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339Nano)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(now))
	return
}
