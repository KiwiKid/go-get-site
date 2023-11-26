package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAI struct {
	client *openai.Client
}

type AIQuestion struct {
	Question           string `json:"question"`
	CorrectAnswer      string `json:"correct_answer"`
	IncorrectAnswerOne string `json:"incorrect_answer_one"`
	IncorrectAnswerTwo string `json:"incorrect_answer_two"`
	QuestionContext    string `json:"question_context"`
}

type QuestionResult struct {
	Questions []Question `json:"questions"`
}

// NewDB creates a new DB instance with GORM
func NewOpenAI() (*OpenAI, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, errors.New("no OPENAI_API_KEY set")
	}

	client := openai.NewClient(openAIKey)

	return &OpenAI{client: client}, nil
}

func (oai *OpenAI) createChatCompletion(content string) (*string, error) {

	fmt.Printf("ChatCompletion request %s", content)

	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	fmt.Printf("ChatCompletion result %s", resp.Choices[0].Message.Content)

	return &resp.Choices[0].Message.Content, nil
}

func (oai *OpenAI) generateQuestions(content string, startIndex uint, endIndex uint) (*QuestionResult, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Generates a multiple-choice question based on the provided context. The question should be relevant to the context and followed by one correct answer and two incorrect answers.",
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "Content for questions: " + substring(content, startIndex, endIndex),
		},
	}

	log.Printf("OPEN AI REQUEST: %v", messages)

	functions := []openai.FunctionDefinition{
		{
			Name:        "question_generator",
			Description: "Generates a multiple-choice question based on the provided context. The question should be relevant to the context and followed by one correct answer and two incorrect answers.",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"question_context": {
						Type:        jsonschema.String,
						Description: "A text passage providing context for the question. The question generated should be directly related to the information in this passage.",
					},
					"question": {
						Type:        jsonschema.String,
						Description: "A question generated based on the provided context. This should be a clear, concise question that can be answered based on the context.",
					},
					"correct_answer": {
						Type:        jsonschema.String,
						Description: "The correct answer to the generated question. This should be factually accurate and directly inferable from the question context.",
					},
					"incorrect_answer_one": {
						Type:        jsonschema.String,
						Description: "The first incorrect answer option for the question. This should be plausible but clearly distinguishable from the correct answer.",
					},
					"incorrect_answer_two": {
						Type:        jsonschema.String,
						Description: "The second incorrect answer option for the question. This should also be plausible but incorrect, providing a reasonable but incorrect alternative.",
					},
				},
				Required: []string{"question_context", "question", "correct_answer", "incorrect_answer_one", "incorrect_answer_two"},
			},
		},
	}
	log.Printf("OPEN AI REQUEST: %v \n %v", messages, functions)

	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  messages,
			Functions: functions,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices found")
	}

	log.Printf("OPEN AI RESPONSE: %v", len(resp.Choices))

	for _, r := range resp.Choices {
		log.Printf("OPEN AI RESPONSE: %v", r)
	}

	// Assuming the response content is in a format that can be unmarshaled into the QuestionResult struct
	var result QuestionResult
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &result, nil
}
