package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient(kubeconfigPath string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) GetPods(namespace string) (*corev1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
}

func (c *Client) GetServices(namespace string) (*corev1.ServiceList, error) {
	return c.clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
}

func (c *Client) GetDeployments(namespace string) (*appsv1.DeploymentList, error) {
	return c.clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
}

func (c *Client) GetNodes() (*corev1.NodeList, error) {
	return c.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
}

func (c *Client) GetNamespaces() (*corev1.NamespaceList, error) {
	return c.clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
}

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
