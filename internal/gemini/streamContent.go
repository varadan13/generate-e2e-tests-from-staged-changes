package gemini

import (
	"context"
	"fmt"
	"generate-e2e/internal/git"
	"log"

	"google.golang.org/genai"
)

func getGenaiClient() (*genai.Client, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  GetApiKey(),
		Backend: genai.BackendGeminiAPI,
	})

	return client, err
}

func StreamContent() {
	ctx := context.Background()

	client, err := getGenaiClient()

	if err != nil {
		log.Fatal(err)
	}

	diff, err := git.GetStagedDiff()

	if err != nil {
		log.Fatal(err)
	}

	systemprompt := fmt.Sprintf(`
					**Prompt:** You are an expert in writing test scenarios for an e-commerce website.Analyze the following Git diff provided in the input and generate comprehensive test scenarios that ensure no existing functionality is broken by the changes.Your test scenarios should adhere to the following guidelines: 1.Focus exclusively on the lines of code that have changed in the diff.2.Simulate realistic user interactions, such as clicking buttons, typing in forms, navigating through the site, and adding items to the cart.3.Consider edge case scenarios that could arise from the changes, including but not limited to: - Input validation errors (e.g., incorrect email formats, empty fields).- Boundary conditions (e.g., maximum character limits for inputs).- Situations with unexpected user behavior (e.g., rapid clicking, navigating away during transactions).- Compatibility with various browsers or devices.- Handling of network failures or slow connections during critical operations.Please format your test scenarios clearly, specifying the user action, the expected outcome, and any necessary setup or preconditions.Ensure that the scenarios cover both positive and negative test cases.
					
					Also add another section called Code Review and review the diff changes brutally.
					
					Git diff: <INSERT_DIFF_BELOW>					
					%s
					`, *diff)

	// prompt := fmt.Sprintf("**Prompt:** You are an AI specialized in generating clear and concise commit messages based on git diffs.Your task is to analyze the provided git diff and summarize the changes in a structured commit message.Follow these guidelines: 1.**Identify the Purpose**: Determine the main purpose of the changes (e.g., bug fix, feature addition, refactoring, documentation update).2.**Summarize Changes**: List the key modifications made, focusing on what files were changed and the type of changes (additions, deletions, modifications).3.**Use Imperative Mood**: Write the commit message in the imperative mood, starting with a verb (e.g., Add, Fix, Update).4.**Limit Length**: Keep the summary line to 50 characters or less, followed by a more detailed explanation if necessary.5.**Include Context**: If there are any related issues or tickets, mention them at the end of the message.6.**Format**: Ensure that the commit message follows conventional commit standards if applicable.**Input Git Diff**: ``` %s ```", *diff)

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: systemprompt},
			},
		},
		// {
		// 	Role: "user",
		// 	Parts: []*genai.Part{
		// 		{Text: prompt},
		// 	},
		// },
	}

	for result, err := range client.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		contents,
		nil,
	) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(result.Candidates[0].Content.Parts[0].Text)
	}

}
