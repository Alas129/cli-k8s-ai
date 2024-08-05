package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/Alas129/cli-k8s-ai/pkg/k8s"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	return &Client{client: openai.NewClient(apiKey)}, nil
}

func (c *Client) ProcessUserInput(input string, k8sClient *k8s.Client) (string, error) {
	ctx := context.Background()

	messages := createChatMessages(input)
	functions := getFunctionDefinitions()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  messages,
			Functions: functions,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %v", err)
	}

	if resp.Choices[0].Message.FunctionCall != nil {
		// Handle function call
		return handleFunctionCall(resp.Choices[0].Message.FunctionCall, k8sClient)
	}

	return resp.Choices[0].Message.Content, nil
}
