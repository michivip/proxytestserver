package webserver

import (
	"net/http"
	"net"
	"strings"
	"log"
)

var proxyHeaders = []string{
	//  proxy variables
	"Via", "X_forwarded", "X_forwarded_for", "X-Forwarded-For", "X-Forwarded",
	// http variables
	"Http_forwarded", "Http-Forwarded", "Http_x_forwarded_for", "Http_client_ip", "Http_via", "Http_proxy_connection", "Http_proxy_connection", "Http-X-Forwarded-For", "Http-Client-Ip",
}

type CheckProxyResponse struct {
	IsProxy                bool            `json:"is_proxy"`
	ProxyHeaders           http.Header     `json:"proxy_headers"`
	SuspectedRealAddresses map[string]*int `json:"suspected_real_address"`
}

func EndpointCheckProxy(writer http.ResponseWriter, req *http.Request) {
	response := &CheckProxyResponse{
		IsProxy:      false,
		ProxyHeaders: http.Header{},
	}
	response.SuspectedRealAddresses = fetchSuspectedIPAddresses(req, response)
	remoteIpString := req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]
	if amount, ok := response.SuspectedRealAddresses[remoteIpString]; ok {
		*amount += 1
	} else {
		defaultValue := 1
		response.SuspectedRealAddresses[remoteIpString] = &defaultValue
	}
	writeJsonResponse(writer, req, response)
}

func fetchSuspectedIPAddresses(req *http.Request, response *CheckProxyResponse) (suspectedAddresses map[string]*int) {
	log.Printf("Proxy check from: %v\n", req.Header)
	suspectedAddresses = make(map[string]*int, 0)
	for _, proxyHeader := range proxyHeaders {
		if value, ok := req.Header[proxyHeader]; ok {
			response.ProxyHeaders[proxyHeader] = req.Header[proxyHeader]
			if !response.IsProxy {
				response.IsProxy = true
			}
			if len(value) == 1 {
				if ip := parseIp(value[0]); ip != "" {
					if value, ok := suspectedAddresses[ip]; ok {
						*value += 1
					} else {
						initialValue := 1
						suspectedAddresses[ip] = &initialValue
					}
				}
			}
		}
	}
	return
}

func parseIp(input string) (ip string) {
	if addresses, err := net.LookupHost(input); err != nil && len(addresses) > 0 {
		return addresses[0]
	} else if host, _, err := net.SplitHostPort(input); err == nil {
		if addresses, err := net.LookupHost(host); err == nil && len(addresses) > 0 {
			return addresses[0]
		} else if ipBytes := net.ParseIP(input); ipBytes != nil {
			return ipBytes.String()
		}
	} else {
		split := strings.Split(input, " ")
		for _, elem := range split {
			if ip = parseIp(elem); ip != "" {
				return ip
			}
		}
	}
	return
}

func EndpointRequestHeaders(writer http.ResponseWriter, req *http.Request) {
	writeJsonResponse(writer, req, req.Header)
}
