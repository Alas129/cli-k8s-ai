package openai

import "github.com/sashabaranov/go-openai"

func createChatMessages(input string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant that answers questions about Kubernetes resources.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}
}
