# CLI app for Kubernetes

A CLI application that interacts with a Kubernetes cluster and provides information through an AI-powered chat interface. The application utilizes OpenAI’s API to process user input and interact with Kubernetes resources.

## Features
- **Chat with Kubernetes**: Start a chat session where you can query your Kubernetes cluster using natural language.
- **Get Cluster Information**: Retrieve detailed information about pods, services, deployments, nodes, and namespaces.
- **Cluster Summary**: Obtain a high-level summary of the cluster’s state.

## Installation
1. Clone the Repository
```bash
git clone https://github.com/your-username/cli-k8s-ai.git
cd cli-k8s-ai
```

2. Build the Application
```
go build -o k8s-cli
```

3. Set Up Environment Variables
Set the OPENAI_API_KEY environment variable to your OpenAI API key.
```bash
export OPENAI_API_KEY=your-openai-api-key
```

## Usage
**Start the Chat Interface**

  Run the application to start the chat session with your Kubernetes cluster.

```bash
./k8s-cli chat
```

  You can now type your queries about the Kubernetes cluster. For example:

```bash
> Give me a summary of the cluster
...
> List all pods in the default namespace
```

## Code Structure

- main.go: Entry point of the application.
- cmd/chat.go: CLI command for interacting with Kubernetes.
- pkg/openai/: Contains OpenAI client logic and handling.
	- client.go: Initializes OpenAI client and processes user input.
	- messages.go: Creates chat messages for OpenAI.
	- functions.go: Defines and handles function calls from OpenAI.
	- formatters.go: Formats Kubernetes resource lists.
- pkg/k8s/: Contains Kubernetes client logic and resource handling.
	- client.go: Initializes Kubernetes client.
	- resources.go: Provides methods to get Kubernetes resources.
	- summary.go: Provides methods to get a cluster summary.