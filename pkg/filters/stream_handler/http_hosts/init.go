package http_hosts

import (
	"encoding/json"
	"strings"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
	"github.com/wzshiming/hostmatcher"
)

func init() {
	registries.StreamHandlerRegistry.Register("http_hosts", NewHttpHosts)
}

type Config struct {
	Hosts []ConfigRoute `json:"hosts"`
}

type ConfigRoute struct {
	Domains []string `json:"domains"`
	Cluster string   `json:"cluster"`
}

func NewHttpHosts(config json.RawMessage) (registries.StreamHandler, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}

	httpHosts := &HttpHosts{
		Hosts: map[string]registries.StreamHandler{},
	}

	for _, host := range conf.Hosts {
		if host.Cluster == "" {
			continue
		}
		route := registries.StreamHandlerInstance.Get(host.Cluster)
		for _, domain := range host.Domains {
			if domain == "" || domain == "*" {
				httpHosts.Default = route
				continue
			}
			if strings.HasPrefix(domain, ".") || strings.Contains(domain, "*") {
				match := hostmatcher.NewMatcher([]string{domain})
				httpHosts.Matches = append(httpHosts.Matches, HttpHostMatches{
					Match:   match,
					Handler: route,
				})
				continue
			}
			httpHosts.Hosts[domain] = route
		}
	}
	if httpHosts.Default == nil {
		httpHosts.Default = registries.StreamHandlerInstance.Get("default")
	}
	return httpHosts, nil
}
