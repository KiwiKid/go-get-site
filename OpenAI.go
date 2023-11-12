package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAI struct {
	client *openai.Client
}

// NewDB creates a new DB instance with GORM
func NewOpenAI() (*OpenAI, error) {
	openAIKey := os.Getenv("OPEN_AI_KEY")
	if openAIKey == "" {
		return nil, errors.New("no OPEN_AI_KEY set")
	}

	client := openai.NewClient(openAIKey)

	return &OpenAI{client: client}, nil
}

func (oai *OpenAI) generateQuestions(content string) {

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	}
	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
			Functions: []openai.FunctionDefinition{
				{
					Name: "get_current_weather",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"question": {
								Type:        jsonschema.String,
								Description: "A question based on the provided content",
							},
							"correct_answer": {
								Type:        jsonschema.String,
								Description: "The city and state, e.g. San Francisco, CA",
							},
							"incorrect_answer_one": {
								Type:        jsonschema.String,
								Description: "The city and state, e.g. San Francisco, CA",
							},
							"incorrect_answer_two": {
								Type:        jsonschema.String,
								Description: "The city and state, e.g. San Francisco, CA",
							},
							"question_context": {
								Type:        jsonschema.String,
								Description: "A quote from the content relating to the answer to the question",
							},
						},
						Required: []string{"question", "question_context", "correct_answer", "incorrect_answer_one", "incorrect_answer_two"},
					},
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
