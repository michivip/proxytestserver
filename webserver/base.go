package webserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"github.com/michivip/proxytestserver/config"
	"strconv"
)

type ReverseProxyHandler struct {
	ReverseProxyHeader string
	RealHandler        http.Handler
}

func (handler *ReverseProxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler.ReverseProxyHeader != "" {
		request.RemoteAddr = request.Header.Get(handler.ReverseProxyHeader)
	}
	userAgentHeader := request.Header["User-Agent"]
	var userAgent string
	if len(userAgentHeader) > 0 {
		userAgent = userAgentHeader[0]
	} else {
		userAgent = "<not set>"
	}
	log.Printf("[%v][%v] Host: %v, User-Agent: %v\n", request.RemoteAddr, request.RequestURI, strconv.Quote(request.Host), strconv.Quote(userAgent))
	handler.RealHandler.ServeHTTP(writer, request)
}

func StartWebserver(config *config.Configuration) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/headers", EndpointRequestHeaders())
	router.HandleFunc("/proxycheck", EndpointCheckProxy(config))
	server := &http.Server{
		Handler: &ReverseProxyHandler{
			ReverseProxyHeader: config.ReverseProxyHeader,
			RealHandler:        router,
		},
		Addr: config.Address,
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
		log.Printf("Could not write headers to %v because of %T: %v\n", req.RemoteAddr, err, err)
	}
}
