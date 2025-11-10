package gemini

import (
	"google.golang.org/genai"
)

func GetConfigWithTool() *genai.GenerateContentConfig {
	tool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_staged_diff",
				Description: "You are a developer. You have made some changes in a repo. Just obtain the staged diff and generate a short meaningful commit message",
				// Parameters: &genai.Schema{
				// 	Type: "object",
				// 	Properties: map[string]*genai.Schema{
				// 		"location": {
				// 			Type:        "string",
				// 			Description: "The city name, e.g. San Francisco",
				// 		},
				// 	}, Required: []string{"location"}},
			}},
	}

	return &genai.GenerateContentConfig{Tools: []*genai.Tool{tool}}
}
