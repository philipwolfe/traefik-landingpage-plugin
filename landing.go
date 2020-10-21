// Package landing_plugin a landing plugin.
package landing_plugin

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	IncludeHost bool `json:"include-host,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// landing a landing plugin.
type landing struct {
	includeHost bool
	next        http.Handler
	name        string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &landing{
		includeHost: config.IncludeHost,
		next:        next,
		name:        name,
	}, nil
}

func (e *landing) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Add("Content-Type", "application/json")
	if e.includeHost {
		resp, _ := json.Marshal(map[string]string{
			"ip":   GetIP(req),
			"path": req.URL.RequestURI(),
		})
		rw.Write(resp)
	} else {
		resp, _ := json.Marshal(map[string]string{
			"ip":   GetIP(req),
			"path": req.URL.RequestURI(),
			"host": GetHost(req),
		})
		rw.Write(resp)
	}

}

//GetIP func
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

//GetHost func
func GetHost(r *http.Request) string {
	if r.URL.IsAbs() {
		host := r.Host
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		return host
	}
	return r.URL.Host
}
