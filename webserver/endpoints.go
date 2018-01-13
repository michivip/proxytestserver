package webserver

import (
	"net/http"
	"strings"
	"regexp"
	"github.com/michivip/proxytestserver/config"
	"errors"
	"log"
)

type CheckProxyResponse struct {
	IsProxy                bool            `json:"is_proxy"`
	ProxyHeaders           http.Header     `json:"proxy_headers"`
	SuspectedRealAddresses map[string]*int `json:"suspected_real_address"`
}

var ErrMaximumHeaderLengthExceeded = errors.New("the maximum header length has been exceeded")

func EndpointCheckProxy(config *config.Configuration) func(http.ResponseWriter, *http.Request) {
	ipRegex, err := regexp.Compile(config.IpRegex)
	if err != nil {
		log.Panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		response := &CheckProxyResponse{
			IsProxy:      false,
			ProxyHeaders: http.Header{},
		}
		var err error

		response.SuspectedRealAddresses, err = fetchSuspectedIPAddresses(request, response, config, ipRegex)
		if err == ErrMaximumHeaderLengthExceeded {
			http.Error(writer, "431 request header fields too large", http.StatusRequestHeaderFieldsTooLarge)
			return
		}
		remoteIpString := request.RemoteAddr[:strings.LastIndex(request.RemoteAddr, ":")]
		if amount, ok := response.SuspectedRealAddresses[remoteIpString]; ok {
			*amount += 1
		} else {
			defaultValue := 1
			response.SuspectedRealAddresses[remoteIpString] = &defaultValue
		}
		writeJsonResponse(writer, request, response)
	}
}

func fetchSuspectedIPAddresses(req *http.Request, response *CheckProxyResponse, config *config.Configuration, ipRegex *regexp.Regexp) (suspectedAddresses map[string]*int, err error) {
	suspectedAddresses = make(map[string]*int, 0)
	for header, headerValue := range req.Header {
		for _, proxyHeader := range config.ProxyHeaders {
			if err = checkProxyRequestHeader(headerValue, proxyHeader, header, response, req, suspectedAddresses, config.MaximumHeaderLength, ipRegex); err != nil {
				return
			}
		}
	}
	return
}

func checkProxyRequestHeader(headerValue []string, proxyHeader string, header string, response *CheckProxyResponse, req *http.Request, suspectedAddresses map[string]*int, maximumHeaderLength int, ipRegex *regexp.Regexp) error {
	for _, valueElem := range headerValue {
		if len(valueElem) > maximumHeaderLength {
			return ErrMaximumHeaderLengthExceeded
		}
		if proxyHeader != header {
			continue
		}
		response.ProxyHeaders[proxyHeader] = req.Header[proxyHeader]
		if !response.IsProxy {
			response.IsProxy = true
		}
		for _, foundIp := range ipRegex.FindAllString(valueElem, -1) {
			if mapIp, ok := suspectedAddresses[foundIp]; ok {
				*mapIp += 1
			} else {
				var initialCount = 1
				suspectedAddresses[foundIp] = &initialCount
			}
		}
	}
	return nil
}

func EndpointRequestHeaders() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, request.Header)
	}
}
