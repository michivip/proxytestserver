# proxytestserver [![Build Status](https://travis-ci.org/michivip/proxytestserver.svg?branch=master)](https://travis-ci.org/michivip/proxytestserver)
A simple webserver to determine a proxy's type written in Golang.
# Installation
## Download executable binary
You can download the binary for your system at the [releases page](https://github.com/michivip/proxytestserver/releases).
## Build your own version
You can get the source by using the built in `go get` command:
```
go get -t github.com/michivip/proxytestserver
```
To build the binary just run the `go build` command:
```
go build
```

# Parameters
- **--config <path>**: The path to your configuration file (default: `config.toml`)
- **--logfile <path>**: The path to your logging file - if empty, no log file is used (default: empty) 

# Configuration
- **address**: Address, the server will bind to.
- **proxy_headers**: A slice of headers which should be checked on /proxycheck (default: most common http/proxy redirection headers) 
- **ip_regex**: The regular expression which searches the headers for ip addresses (default: a regex which searches for ipv4/ipv6).
- **maximum_header_length**: The maximum amount of characters in a header value (to prevent attacks) (default: `128`)
- **reverse_proxy_header**: The header name for a reverse proxy which contains the real ip address in the format `<ADDR>:<PORT>` - if empty, no real address will be fetched (default: empty)

# Contributing
If you want to contribute, just open an issue. Then your issue will be discussed.

# Used libraries
- [TOML parser by BurntSushi](https://github.com/BurntSushi/toml)