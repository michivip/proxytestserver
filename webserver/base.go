package webserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"github.com/michivip/proxytestserver/config"
)

func StartWebserver(config *config.Configuration) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/headers", EndpointRequestHeaders())
	router.HandleFunc("/proxycheck", EndpointCheckProxy(config))
	server := &http.Server{
		Handler: router,
		Addr:    config.Address,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
	return server
}

func writeJsonResponse(writer http.ResponseWriter, req *http.Request, value interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		http.Error(writer, "500 Internal error", http.StatusInternalServerError)
		log.Printf("Could not write headers to %v because of %T: %v", req.RemoteAddr, err, err)
	}
}
