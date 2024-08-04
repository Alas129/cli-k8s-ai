/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Alas129/cli-k8s-ai/pkg/k8s"
	"github.com/Alas129/cli-k8s-ai/pkg/openai"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with your Kubernetes cluster",
	Long:  `Start a chat session with your Kubernetes cluster`,
	Run:   runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)
}

func runChat(cmd *cobra.Command, args []string) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	k8sClient, err := k8s.NewClient(kubeconfig)
	if err != nil {
		fmt.Printf("Error initializing Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	openAIClient, err := openai.NewClient()
	if err != nil {
		fmt.Printf("Error initializing OpenAI client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Welcome to the Kubernetes ChatBot. Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			break
		}

		response, err := openAIClient.ProcessUserInput(input, k8sClient)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(response)
		}
	}
}
