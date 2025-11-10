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
					You are an expert in writing E2E tests using Playwright and Jest.

					Analyze the following Git diff and generate an appropriate end-to-end test
					that ensures the changed behavior is correct.

					Guidelines:
					- Focus only on what changed.
					- Simulate realistic user interactions (e.g., clicking, typing, navigating).
					- Include setup, teardown, and assertions.
					- Use existing project conventions if visible.
					- Output runnable test code and a short summary.

					Git diff:
					<INSERT_DIFF_HERE>
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
