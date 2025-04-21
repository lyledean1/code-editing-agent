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

// MockClient is a mock implementation of the Anthropic client
type MockClient struct {
	mock.Mock
	anthropic.Client
}

// MockMessagesService mocks the Messages service of the Anthropic client
type MockMessagesService struct {
	mock.Mock
}

// New mocks the New method of the Messages service
func (m *MockMessagesService) New(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*anthropic.Message), args.Error(1)
}

func TestNewAgent(t *testing.T) {
	// Arrange
	mockClient := &anthropic.Client{}
	getUserMessage := func() (string, bool) { return "", true }
	mockTools := []tools.ToolDefinition{tools.ReadFileDefinition}

	// Act
	agent := NewAgent(mockClient, getUserMessage, mockTools)

	// Assert
	assert.NotNil(t, agent)
	assert.Equal(t, mockClient, agent.client)
	assert.Equal(t, mockTools, agent.tools)
	// Note: Can't directly compare functions, but at least check it's not nil
	assert.NotNil(t, agent.getUserMessage)
}

func TestExecuteTool(t *testing.T) {
	// Arrange
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

	// Act
	result := agent.executeTool("tool-123", "mock_tool", json.RawMessage(`{}`))

	// Assert
	// Since the result is a ContentBlockParamUnion, we need to check properties without using GetToolResult()
	// Let's verify the basic equality of what we expected
	expected := anthropic.NewToolResultBlock("tool-123", "mock result", false)
	assert.Equal(t, expected, result)
}

func TestExecuteToolNotFound(t *testing.T) {
	// Arrange
	mockClient := &anthropic.Client{}
	getUserMessage := func() (string, bool) { return "", true }
	agent := NewAgent(mockClient, getUserMessage, []tools.ToolDefinition{})

	// Act
	result := agent.executeTool("tool-123", "non_existent_tool", json.RawMessage(`{}`))

	// Assert
	// Since the result is a ContentBlockParamUnion, we need to check properties without using GetToolResult()
	// Let's verify the basic equality of what we expected
	expected := anthropic.NewToolResultBlock("tool-123", "tool not found", true)
	assert.Equal(t, expected, result)
}
