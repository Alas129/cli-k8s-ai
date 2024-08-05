package k8s

import (
	"context"

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
