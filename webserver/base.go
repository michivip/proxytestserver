package webserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/michivip/proxytestserver/config"
	"log"
	"net/http"
	"strconv"
)

type reverseProxyHandler struct {
	ReverseProxyHeader string
	RealHandler        http.Handler
}

func (handler *reverseProxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

// this method asynchronous starts the web server with the provided configuration values. The server runs in the background but with the returned pointer, interaction with the server is given.
func StartWebserver(config *config.Configuration) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/headers", endpointRequestHeaders())
	router.HandleFunc("/proxycheck", endpointCheckProxy(config))
	server := &http.Server{
		Handler: &reverseProxyHandler{
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
