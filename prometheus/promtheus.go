package prometheus

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shaxiaozz/prometheus-mcp-server/utils"
	"github.com/wonderivan/logger"
	"time"
)

const (
	ListMetrics       = "prometheus_list_metrics"
	GetTargets        = "prometheus_get_targets"
	GetMetricMetadata = "prometheus_get_metric_metadata"
	ExecuteQuery      = "prometheus_execute_query"
	ExecuteRangeQuery = "prometheus_execute_range_query"
	ExecuteLastQuery  = "prometheus_execute_last_query"
)

var (
	ListMetricsTool       = mcp.NewTool(ListMetrics, mcp.WithDescription("List all available metrics in Prometheus"))
	GetTargetsTool        = mcp.NewTool(GetTargets, mcp.WithDescription("Get information about all scrape targets"))
	GetMetricMetadataTool = mcp.NewTool(GetMetricMetadata,
		mcp.WithDescription("Get metadata for a specific metric"),
		mcp.WithString("metric",
			mcp.Description("The name of the metric to retrieve metadata for"),
			mcp.Required()),
	)
	ExecuteQueryTool = mcp.NewTool(ExecuteQuery,
		mcp.WithDescription("Execute a PromQL instant query against Prometheus"),
		mcp.WithString("query",
			mcp.Description("Prometheus expression query string."),
			mcp.Required()),
		mcp.WithString("time",
			mcp.Description("Optional RFC3339 or Unix timestamp")),
	)
	ExecuteRangeQueryTool = mcp.NewTool(ExecuteRangeQuery,
		mcp.WithDescription("Execute a PromQL range query with start time, end time, and step interval. (The maximum data size cannot exceed 30 minutes (within the start and end time))"),
		mcp.WithString("query",
			mcp.Description("Prometheus expression query string."),
			mcp.Required()),
		mcp.WithString("start",
			mcp.Description("Start time as RFC3339 or Unix timestamp"),
			mcp.Required()),
		mcp.WithString("end",
			mcp.Description("End time as RFC3339 or Unix timestamp"),
			mcp.Required()),
		mcp.WithString("step",
			mcp.Description("Query resolution step width in duration format or float number of seconds. (e.g., '15s', '1m', '1h')"),
			mcp.DefaultString("15s")),
		mcp.WithString("timeout",
			mcp.Description("Evaluation timeout."),
			mcp.DefaultString("15s")),
	)
	ExecuteLastQueryTool = mcp.NewTool(ExecuteLastQuery,
		mcp.WithDescription("Executes a PromQL range query with a recent time and a step interval. The recent time is in minutes. (e.g., '5', '10', '30') (The maximum data size cannot exceed 30 minutes)"),
		mcp.WithString("query",
			mcp.Description("Prometheus expression query string."),
			mcp.Required()),
		mcp.WithString("last_minute",
			mcp.Description("The recent time is in minutes. (e.g., '5m', '10m', '30m')"),
			mcp.Required()),
		mcp.WithString("step",
			mcp.Description("Query resolution step width in duration format or float number of seconds. (e.g., '15s', '1m', '1h')"),
			mcp.DefaultString("15s")),
		mcp.WithString("timeout",
			mcp.Description("Evaluation timeout."),
			mcp.DefaultString("15s")),
	)
)

func ListMetricsToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call PrometheusListMetricsToolHandle:  ")
	data, err := utils.MakePrometheusRequest("label/__name__/values", nil)
	if err != nil {
		errorMsg := "Call PrometheusListMetricsToolHandle Failure: " + err.Error()
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call PrometheusListMetricsToolHandle data: ")
	//fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}

func GetTargetsToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call GetTargetsToolHandle:  ")
	data, err := utils.MakePrometheusRequest("targets", nil)
	if err != nil {
		errorMsg := "Call GetTargetsToolHandle Failure: " + err.Error()
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call GetTargetsToolHandle data: ")
	//fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}

func GetMetricMetadataToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call GetMetricMetadataToolHandle:  ")
	errorMsg := "Call GetMetricMetadataToolHandle Failure: "

	metric, ok := request.Params.Arguments["metric"].(string)
	if !ok || metric == "" {
		errorMsg = errorMsg + "metric must be a string or empty"
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}

	params := map[string]string{"metric": metric}
	data, err := utils.MakePrometheusRequest("metadata", params)
	if err != nil {
		errorMsg = errorMsg + err.Error()
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call GetMetricMetadataToolHandle " + metric + " data: ")
	fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}

func ExecuteQueryToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call ExecuteQueryToolHandle:  ")
	errorMsg := "Call ExecuteQueryToolHandle Failure: "

	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		errorMsg = errorMsg + "query must be a string or empty"
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}

	params := map[string]string{"query": query}
	data, err := utils.MakePrometheusRequest("query", params)
	if err != nil {
		errorMsg = errorMsg + err.Error()
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call ExecuteQueryToolHandle " + query + " data: ")
	fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}

func ExecuteRangeQueryToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call ExecuteRangeQueryToolHandle:  ")
	errorMsg := "Call ExecuteRangeQueryToolHandle Failure: "

	params := map[string]string{}
	for k, v := range request.Params.Arguments {
		str, ok := v.(string)
		if !ok || str == "" {
			errorMsg = errorMsg + k + " must be a string or empty"
			logger.Info(errorMsg)
			return nil, errors.New(errorMsg)
		}
		params[k] = str
	}

	data, err := utils.MakePrometheusRequest("query_range", params)
	if err != nil {
		errorMsg = errorMsg + err.Error()
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call ExecuteRangeQueryToolHandle " + fmt.Sprintf("%#v", params) + " data: ")
	fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}

func ExecuteLastQueryToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call ExecuteLastQueryToolHandle:  ")
	errorMsg := "Call ExecuteLastQueryToolHandle Failure: "

	params := map[string]string{}
	lastMinute := "10m"
	for k, v := range request.Params.Arguments {
		str, ok := v.(string)
		if !ok || str == "" {
			errorMsg = errorMsg + k + " must be a string or empty"
			logger.Info(errorMsg)
			return nil, errors.New(errorMsg)
		}
		if k == "last_minute" {
			lastMinute = str
			continue
		}
		params[k] = str
	}

	// start and end times
	end := time.Now().UTC()
	lastDuration, _ := time.ParseDuration(lastMinute)
	start := end.Add(-lastDuration)
	endStr := end.Format(time.RFC3339)
	startStr := start.Format(time.RFC3339)
	params["start"] = startStr
	params["end"] = endStr

	data, err := utils.MakePrometheusRequest("query_range", params)
	if err != nil {
		errorMsg = errorMsg + err.Error()
		logger.Info(errorMsg)
		return nil, errors.New(errorMsg)
	}
	logger.Info("Call ExecuteLastQueryToolHandle " + fmt.Sprintf("%#v", params) + " data: ")
	fmt.Println(string(data))
	return mcp.NewToolResultText(string(data)), nil
}
