package k8s

import (
	"fmt"
)

func (c *Client) GetClusterSummary() (string, error) {
	namespaces, err := c.GetNamespaces()
	if err != nil {
		return "", err
	}

	nodes, err := c.GetNodes()
	if err != nil {
		return "", err
	}

	var summary string

	summary += fmt.Sprintf("Cluster has %d namespaces:\n", len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		summary += fmt.Sprintf("- %s\n", namespace.Name)
	}

	summary += fmt.Sprintf("\nCluster has %d nodes:\n", len(nodes.Items))
	for _, node := range nodes.Items {
		summary += fmt.Sprintf("- %s\n", node.Name)
	}

	for _, namespace := range namespaces.Items {
		pods, err := c.GetPods(namespace.Name)
		if err != nil {
			return "", err
		}
		services, err := c.GetServices(namespace.Name)
		if err != nil {
			return "", err
		}
		deployments, err := c.GetDeployments(namespace.Name)
		if err != nil {
			return "", err
		}

		summary += fmt.Sprintf("\nNamespace '%s' has %d pods, %d services, and %d deployments.\n",
			namespace.Name, len(pods.Items), len(services.Items), len(deployments.Items))
	}

	return summary, nil
}
