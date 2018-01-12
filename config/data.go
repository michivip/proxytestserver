package config

type Loader interface {
	Load() (*Configuration, error)
	Save(*Configuration) (error)
}

type Configuration struct {
	Address      string   `toml:"address"`
	ProxyHeaders []string `toml:"proxy_headers"`
	IpRegex      string   `toml:"ip_regex"`
}

var DefaultConfiguration = &Configuration{
	Address: "localhost:8070",
	ProxyHeaders: []string{
		//  proxy variables
		"Via", "X_forwarded", "X_forwarded_for", "X-Forwarded-For", "X-Forwarded",
		// http variables
		"Http_forwarded", "Http-Forwarded", "Http_x_forwarded_for", "Http_client_ip", "Http_via", "Http_proxy_connection", "Http_proxy_connection", "Http-X-Forwarded-For", "Http-Client-Ip",
	},
	IpRegex: `(?:(?:(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))|(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:)(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)(?:.(?:25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(?:%.+)?s*)`,
}
