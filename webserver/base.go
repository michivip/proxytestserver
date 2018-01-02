package webserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func StartWebserver(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/headers", EndpointRequestHeaders)
	server := &http.Server{
		Handler: router,
		Addr:    addr,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
	return server
}
