package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
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
		{
			Name: "get_services",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "The namespace to list services from",
					},
				},
				"required": []string{"namespace"},
			},
		},
		{
			Name: "get_deployments",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "The namespace to list deployments from",
					},
				},
				"required": []string{"namespace"},
			},
		},
		{
			Name: "get_nodes",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name: "get_namespaces",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name: "get_cluster_summary",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
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
		case "get_services":
			var args struct {
				Namespace string `json:"namespace"`
			}
			if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &args); err != nil {
				return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
			}
			services, err := k8sClient.GetServices(args.Namespace)
			if err != nil {
				return "", fmt.Errorf("error getting services: %v", err)
			}
			return formatServiceList(services), nil
		case "get_deployments":
			var args struct {
				Namespace string `json:"namespace"`
			}
			if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &args); err != nil {
				return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
			}
			deployments, err := k8sClient.GetDeployments(args.Namespace)
			if err != nil {
				return "", fmt.Errorf("error getting deployments: %v", err)
			}
			return formatDeploymentList(deployments), nil
		case "get_nodes":
			nodes, err := k8sClient.GetNodes()
			if err != nil {
				return "", fmt.Errorf("error getting nodes: %v", err)
			}
			return formatNodeList(nodes), nil
		case "get_namespaces":
			namespaces, err := k8sClient.GetNamespaces()
			if err != nil {
				return "", fmt.Errorf("error getting namespaces: %v", err)
			}
			return formatNamespaceList(namespaces), nil
		case "get_cluster_summary":
			summary, err := k8sClient.GetClusterSummary()
			if err != nil {
				return "", fmt.Errorf("error getting cluster summary: %v", err)
				}
			return summary, nil
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

func formatServiceList(services *corev1.ServiceList) string {
	result := "Services:\n"
	for _, service := range services.Items {
		result += fmt.Sprintf("- %s (Type: %s)\n", service.Name, service.Spec.Type)
	}
	return result
}

func formatDeploymentList(deployments *appsv1.DeploymentList) string {
	result := "Deployments:\n"
	for _, deployment := range deployments.Items {
		result += fmt.Sprintf("- %s (Replicas: %d)\n", deployment.Name, *deployment.Spec.Replicas)
	}
	return result
}

func formatNodeList(nodes *corev1.NodeList) string {
	result := "Nodes:\n"
	for _, node := range nodes.Items {
		result += fmt.Sprintf("- %s (Status: %s)\n", node.Name, node.Status.Phase)
	}
	return result
}

func formatNamespaceList(namespaces *corev1.NamespaceList) string {
	result := "Namespaces:\n"
	for _, namespace := range namespaces.Items {
		result += fmt.Sprintf("- %s\n", namespace.Name)
	}
	return result
}
