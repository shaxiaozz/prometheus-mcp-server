package utils

import (
	"encoding/json"
	"fmt"
	"github.com/shaxiaozz/prometheus-mcp-server/config"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PrometheusResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error,omitempty"`
}

func getPrometheusAuth() (map[string]string, *url.Userinfo) {
	// Get authentication for Prometheus based on provided credentials.
	if config.PrometheusToken != "" {
		return map[string]string{
			"Authorization": "Bearer " + config.PrometheusToken,
		}, nil
	} else if config.PrometheusUserName != "" && config.PrometheusPassword != "" {
		return nil, url.UserPassword(config.PrometheusUserName, config.PrometheusPassword)
	}

	return nil, nil
}

func MakePrometheusRequest(endpoint string, params map[string]string) ([]byte, error) {
	baseURL := strings.TrimRight(config.PrometheusUrl, "/") + "/api/v1/" + endpoint
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// 添加查询参数
	query := reqURL.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	reqURL.RawQuery = query.Encode()

	// 创建请求
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	headers, basicAuth := getPrometheusAuth()
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if basicAuth != nil {
		password, _ := basicAuth.Password()
		req.SetBasicAuth(basicAuth.Username(), password)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, body)
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("Prometheus API error: %s", result.Error)
	}

	return result.Data, nil
}
