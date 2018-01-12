package webserver

import (
	"net/http"
	"strings"
	"log"
	"regexp"
	"fmt"
)

var proxyHeaders = []string{
	//  proxy variables
	"Via", "X_forwarded", "X_forwarded_for", "X-Forwarded-For", "X-Forwarded",
	// http variables
	"Http_forwarded", "Http-Forwarded", "Http_x_forwarded_for", "Http_client_ip", "Http_via", "Http_proxy_connection", "Http_proxy_connection", "Http-X-Forwarded-For", "Http-Client-Ip",
}

var ipRegex = regexp.MustCompile(`(?:(?:(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:)(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(?:%.+)?s*)`)

type CheckProxyResponse struct {
	IsProxy                bool            `json:"is_proxy"`
	ProxyHeaders           http.Header     `json:"proxy_headers"`
	SuspectedRealAddresses map[string]*int `json:"suspected_real_address"`
}

var ErrMaximumHeaderLengthExceeded = fmt.Errorf("the maximum header length (%d) has been exceeded", MaximumHeaderLength)

func EndpointCheckProxy(writer http.ResponseWriter, req *http.Request) {
	response := &CheckProxyResponse{
		IsProxy:      false,
		ProxyHeaders: http.Header{},
	}
	var err error
	response.SuspectedRealAddresses, err = fetchSuspectedIPAddresses(req, response)
	if err == ErrMaximumHeaderLengthExceeded {
		http.Error(writer, "431 request header fields too large", http.StatusRequestHeaderFieldsTooLarge)
		return
	}
	remoteIpString := req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]
	if amount, ok := response.SuspectedRealAddresses[remoteIpString]; ok {
		*amount += 1
	} else {
		defaultValue := 1
		response.SuspectedRealAddresses[remoteIpString] = &defaultValue
	}
	writeJsonResponse(writer, req, response)
}

const (
	MaximumHeaderLength = 128
)

func fetchSuspectedIPAddresses(req *http.Request, response *CheckProxyResponse) (suspectedAddresses map[string]*int, err error) {
	log.Printf("Proxy check from: %v\n", req.Header)
	suspectedAddresses = make(map[string]*int, 0)
	for header, headerValue := range req.Header {
		for _, proxyHeader := range proxyHeaders {
			if err = checkProxyRequestHeader(headerValue, proxyHeader, header, response, req, suspectedAddresses); err != nil {
				return
			}
		}
	}
	return
}

func checkProxyRequestHeader(headerValue []string, proxyHeader string, header string, response *CheckProxyResponse, req *http.Request, suspectedAddresses map[string]*int) error {
	for _, valueElem := range headerValue {
		if len(valueElem) > MaximumHeaderLength {
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

func EndpointRequestHeaders(writer http.ResponseWriter, req *http.Request) {
	writeJsonResponse(writer, req, req.Header)
}
