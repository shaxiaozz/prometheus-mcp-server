package prometheus

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/wonderivan/logger"
)

const (
	HelloWorld = "hello_world"
)

var HelloWorldTool = mcp.NewTool(HelloWorld, mcp.WithDescription(
	"Say hello to someone"))

func HelloWorldToolHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Call HelloWorldToolHandle")
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
