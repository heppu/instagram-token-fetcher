package igtoken

import (
	"fmt"
	"net"
	"net/http"

	"github.com/braintree/manners"
)

func startServer(port, path string) (func() bool, error) {
	mux := http.NewServeMux()

	server := manners.NewWithServer(&http.Server{
		Addr:    ":" + port,
		Handler: mux,
	})

	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "")
	})

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	go func() {
		server.Serve(listener)
	}()

	return server.BlockingClose, err
}
