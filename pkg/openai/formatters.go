package openai

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

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
