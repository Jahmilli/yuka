package client

// DefaultBaseURL is the default value for the Options.BaseURL
var DefaultBaseURL = "http://localhost:8080"

// DefaultMaxConnections is the default value for Options.MaxConnections
var DefaultMaxConnections = 10

// Options for connecting to a localtunnel server
type Options struct {
	Subdomain      string
	BaseURL        string
	MaxConnections int
}

func (o *Options) setDefaults() {
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}
	if o.MaxConnections == 0 {
		o.MaxConnections = DefaultMaxConnections
	}
}
