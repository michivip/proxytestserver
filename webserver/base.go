package webserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"github.com/michivip/proxytestserver/config"
)

type ReverseProxyHandler struct {
	ReverseProxyHeader string
	RealHandler        http.Handler
}

func (handler *ReverseProxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	request.RemoteAddr = request.Header.Get(handler.ReverseProxyHeader)
	handler.RealHandler.ServeHTTP(writer, request)
}

func StartWebserver(config *config.Configuration) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/headers", EndpointRequestHeaders())
	router.HandleFunc("/proxycheck", EndpointCheckProxy(config))
	var handler http.Handler
	if config.ReverseProxyHeader != "" {
		handler = &ReverseProxyHandler{
			ReverseProxyHeader: config.ReverseProxyHeader,
			RealHandler:        router,
		}
	} else {
		handler = router
	}
	server := &http.Server{
		Handler: handler,
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
