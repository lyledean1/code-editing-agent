package agent

import (
	"agent/tools"
	"context"
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
	anthropic.Client
}

type MockMessagesService struct {
	mock.Mock
}

func (m *MockMessagesService) New(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*anthropic.Message), args.Error(1)
}

func TestNewAgent(t *testing.T) {
	mockClient := &anthropic.Client{}
	getUserMessage := func() (string, bool) { return "", true }
	mockTools := []tools.ToolDefinition{tools.ReadFileDefinition}

	agent := NewAgent(mockClient, getUserMessage, mockTools)
	assert.NotNil(t, agent)
	assert.Equal(t, mockClient, agent.client)
	assert.Equal(t, mockTools, agent.tools)
	assert.NotNil(t, agent.getUserMessage)
}

func TestExecuteTool(t *testing.T) {
	mockClient := &anthropic.Client{}
	getUserMessage := func() (string, bool) { return "", true }

	// Create a simple mock tool for testing
	mockToolDef := tools.ToolDefinition{
		Name:        "mock_tool",
		Description: "Mock tool for testing",
		InputSchema: anthropic.ToolInputSchemaParam{},
		Function: func(input json.RawMessage) (string, error) {
			return "mock result", nil
		},
	}

	agent := NewAgent(mockClient, getUserMessage, []tools.ToolDefinition{mockToolDef})

	result := agent.executeTool("tool-123", "mock_tool", json.RawMessage(`{}`))
	expected := anthropic.NewToolResultBlock("tool-123", "mock result", false)
	assert.Equal(t, expected, result)
}

func TestExecuteToolNotFound(t *testing.T) {
	mockClient := &anthropic.Client{}
	getUserMessage := func() (string, bool) { return "", true }
	agent := NewAgent(mockClient, getUserMessage, []tools.ToolDefinition{})
	result := agent.executeTool("tool-123", "non_existent_tool", json.RawMessage(`{}`))
	expected := anthropic.NewToolResultBlock("tool-123", "tool not found", true)
	assert.Equal(t, expected, result)
}
