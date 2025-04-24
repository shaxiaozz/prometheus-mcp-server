package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shaxiaozz/prometheus-mcp-server/config"
	"github.com/shaxiaozz/prometheus-mcp-server/prometheus"
	"github.com/shaxiaozz/prometheus-mcp-server/utils"
	"github.com/wonderivan/logger"
	"log"
)

func init() {
	config.Config.InitConfig()
}

var (
	Version = utils.Version
)

func newMCPServer() *server.MCPServer {
	return server.NewMCPServer(
		"github.com/shaxiaozz/prometheus-mcp-server",
		Version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)
}

func addTools(s *server.MCPServer) {
	s.AddTool(prometheus.HelloWorldTool, prometheus.HelloWorldToolHandle)               // Say hello to someone
	s.AddTool(prometheus.ListMetricsTool, prometheus.ListMetricsToolHandle)             // List all available metrics in Prometheus
	s.AddTool(prometheus.GetTargetsTool, prometheus.GetTargetsToolHandle)               // Get information about all scrape targets
	s.AddTool(prometheus.GetMetricMetadataTool, prometheus.GetMetricMetadataToolHandle) // Get metadata for a specific metric
	s.AddTool(prometheus.ExecuteQueryTool, prometheus.ExecuteQueryToolHandle)           // Execute a PromQL instant query against Prometheus
	s.AddTool(prometheus.ExecuteRangeQueryTool, prometheus.ExecuteRangeQueryToolHandle) // Execute a PromQL range query with start time, end time, and step interval.
	s.AddTool(prometheus.ExecuteLastQueryTool, prometheus.ExecuteLastQueryToolHandle)   // Executes a PromQL range query with a recent time and a step interval. The recent time is in minutes
}

func runServer(transport string) error {
	mcpServer := newMCPServer()
	addTools(mcpServer)
	ipAddr := utils.GetIPAddr()
	address := ipAddr + ":8000"

	if transport == "sse" {
		logger.Info("SSE server listening on http://" + address)
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://"+address))
		if err := sseServer.Start(address); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		logger.Info("Run Stdio server")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
	return nil
}

func main() {
	if err := runServer(config.McpServerTransport); err != nil {
		fmt.Printf("server run error: %v\n", err)
		panic(err)
	}
}
