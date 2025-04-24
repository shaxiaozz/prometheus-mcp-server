package config

import (
	"log"
	"os"
)

var Config config
var (
	PrometheusUrl, PrometheusUserName, PrometheusPassword, PrometheusToken, McpServerTransport string
)

type config struct {
}

func (c *config) InitConfig() {
	if os.Getenv("PROMETHEUS_URL") == "" {
		log.Fatalf("PROMETHEUS_URL env variable not set")
	} else {
		PrometheusUrl = os.Getenv("PROMETHEUS_URL")
	}

	if os.Getenv("PROMETHEUS_USERNAME") != "" {
		PrometheusUserName = os.Getenv("PROMETHEUS_USERNAME")
	}

	if os.Getenv("PROMETHEUS_PASSWORD") != "" {
		PrometheusPassword = os.Getenv("PROMETHEUS_PASSWORD")
	}

	if os.Getenv("PROMETHEUS_TOKEN") != "" {
		PrometheusToken = os.Getenv("PROMETHEUS_TOKEN")
	}

	if os.Getenv("MCP_SERVER_TRANSPORT") == "" {
		McpServerTransport = "sse"
	} else {
		McpServerTransport = os.Getenv("MCP_SERVER_TRANSPORT")
	}
}
