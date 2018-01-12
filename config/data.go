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
