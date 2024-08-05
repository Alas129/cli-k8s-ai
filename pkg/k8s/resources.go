package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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