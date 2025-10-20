package main

import (
	"context"
	"net/http"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

const (
	// We use a local LLM running in Ollama to ask a question: https://github.com/ollama/ollama
	providerUrl = "http://localhost:11434/v1"
	llmModel    = "qwen3:4b"

// using openai
// llmModel = "gpt-5-nano"
)

// There are many different ways to provide the context to the LLM.
// You can pass each context as user message, or the list as one user message,
// or pass it in the system prompt. The system prompt itself also has a big impact
// on how well the LLM handles the context, especially for LLMs with < 7B parameters.
// The prompt engineering is up to you, it's out of scope for the vector database.
var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
You are a meeting assistant with strong technical skills, like an IT Architect and Requirement Engineer and Devops Expert.

Answer the question in a very concise manner. Use an unbiased and journalistic tone. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
Answer the question solely based on the provided search results from the meeting notes, just say that you don't know. Don't make anything up.

Anything between the following 'context' XML blocks is retrieved from the meeting notes, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}

Don't mention the meeting notes, context or search results in your answer.
`))

func askLLM(ctx context.Context, contexts []string, question string) string {
	// Using OpenAI
	// openAIClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Using Ollama
	// We can use the OpenAI client because Ollama is compatible with OpenAI's API.
	openAIClient := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    providerUrl,
		HTTPClient: http.DefaultClient,
	})
	sb := &strings.Builder{}
	err := systemPromptTpl.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: sb.String(),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "Question: " + question,
		},
	}
	res, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    llmModel,
		Messages: messages,
	})
	if err != nil {
		panic(err)
	}
	reply := res.Choices[0].Message.Content
	reply = strings.TrimSpace(reply)

	return reply
}
