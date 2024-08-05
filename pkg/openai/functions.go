package openai

import (
	"encoding/json"
	"fmt"

	"github.com/Alas129/cli-k8s-ai/pkg/k8s"
	"github.com/sashabaranov/go-openai"
)

func getFunctionDefinitions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
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
}

func handleFunctionCall(functionCall *openai.FunctionCall, k8sClient *k8s.Client) (string, error) {
	switch functionCall.Name {
	case "get_pods":
		return handleGetPods([]byte(functionCall.Arguments), k8sClient)
	case "get_services":
		return handleGetServices([]byte(functionCall.Arguments), k8sClient)
	case "get_deployments":
		return handleGetDeployments([]byte(functionCall.Arguments), k8sClient)
	case "get_nodes":
		return handleGetNodes(k8sClient)
	case "get_namespaces":
		return handleGetNamespaces(k8sClient)
	case "get_cluster_summary":
		return handleGetClusterSummary(k8sClient)
	default:
		return "", fmt.Errorf("unknown function call: %s", functionCall.Name)
	}
}

func handleGetPods(arguments json.RawMessage, k8sClient *k8s.Client) (string, error) {
	var args struct {
		Namespace string `json:"namespace"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
	}
	pods, err := k8sClient.GetPods(args.Namespace)
	if err != nil {
		return "", fmt.Errorf("error getting pods: %v", err)
	}
	return formatPodList(pods), nil
}

func handleGetServices(arguments json.RawMessage, k8sClient *k8s.Client) (string, error) {
	var args struct {
		Namespace string `json:"namespace"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
	}
	services, err := k8sClient.GetServices(args.Namespace)
	if err != nil {
		return "", fmt.Errorf("error getting services: %v", err)
	}
	return formatServiceList(services), nil
}

func handleGetDeployments(arguments json.RawMessage, k8sClient *k8s.Client) (string, error) {
	var args struct {
		Namespace string `json:"namespace"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("error unmarshaling function arguments: %v", err)
	}
	deployments, err := k8sClient.GetDeployments(args.Namespace)
	if err != nil {
		return "", fmt.Errorf("error getting deployments: %v", err)
	}
	return formatDeploymentList(deployments), nil
}

func handleGetNodes(k8sClient *k8s.Client) (string, error) {
	nodes, err := k8sClient.GetNodes()
	if err != nil {
		return "", fmt.Errorf("error getting nodes: %v", err)
	}
	return formatNodeList(nodes), nil
}

func handleGetNamespaces(k8sClient *k8s.Client) (string, error) {
	namespaces, err := k8sClient.GetNamespaces()
	if err != nil {
		return "", fmt.Errorf("error getting namespaces: %v", err)
	}
	return formatNamespaceList(namespaces), nil
}

func handleGetClusterSummary(k8sClient *k8s.Client) (string, error) {
	summary, err := k8sClient.GetClusterSummary()
	if err != nil {
		return "", fmt.Errorf("error getting cluster summary: %v", err)
	}
	return summary, nil
}
