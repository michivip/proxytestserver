package webserver

import (
	"net/http"
	"encoding/json"
	"log"
)

func EndpointRequestHeaders(writer http.ResponseWriter, req *http.Request) {
	if err := json.NewEncoder(writer).Encode(req.Header); err != nil {
		http.Error(writer, "500 Internal error", http.StatusInternalServerError)
		log.Printf("Could not write headers to %v because of %T: %v", req.RemoteAddr, err, err)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	}
}
