package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return &Client{
		client: openai.NewClient(apiKey),
	}, nil
}

func (c *Client) ProcessUserInput(input string, k8sClient *k8s.Client) (string, error) {
	ctx := context.Background()

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant that answers questions about Kubernetes resources.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	functions := []openai.FunctionDefinition{
		{
			Name: "get_pods",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "The namespace to list pods from",
					},
				},
				"required": []string{"namespace"},
			},
		},
	}

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
		switch resp.Choices[0].Message.FunctionCall.Name {
		case "get_pods":
			var args struct {
				Namespace string `json:"namespace"`
			}
			if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &args); err != nil {
				return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
			}
			pods, err := k8sClient.GetPods(args.Namespace)
			if err != nil {
				return "", fmt.Errorf("error getting pods: %v", err)
			}
			return formatPodList(pods), nil
		}
	}

	return resp.Choices[0].Message.Content, nil
}

func formatPodList(pods *corev1.PodList) string {
	result := "Pods:\n"
	for _, pod := range pods.Items {
		result += fmt.Sprintf("- %s (Status: %s)\n", pod.Name, pod.Status.Phase)
	}
	return result
}
