package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

var (
	port = flag.String("port", "80", "listen port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		Handler:     bcophttp.NewHandler(http.DefaultServeMux, propagation.EnvID{}),
		ConnContext: bcophttp.ConnContext,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	bln := bcopnet.NewListener(ln)
	fmt.Println("serve http")
	log.Fatal(server.Serve(bln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339Nano)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(now))
	return
}
