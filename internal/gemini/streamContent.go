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
**Prompt:** You are a senior JavaScript engineer and debugging specialist with deep expertise in frontend and backend JavaScript runtimes.

Analyze the following Git diff and identify **all potential runtime JavaScript errors** that could be introduced by these changes.

Your analysis must strictly follow these rules:

1. **Focus ONLY on the changed lines in the diff** — do not speculate about unchanged code.
2. Identify issues that could cause **runtime failures**, including but not limited to:
   - Undefined or null access (e.g., reading properties of undefined/null)
   - Incorrect assumptions about data shape or types
   - Missing optional chaining or defensive checks
   - Async/await misuse (unhandled promises, missing awaits, race conditions)
   - Event handler issues (stale closures, incorrect bindings)
   - State management errors (stale state, improper mutations)
   - Browser-only APIs used in non-browser environments (SSR, Node)
   - Incorrect imports, renamed variables, or shadowed variables
   - Breaking changes caused by refactors (signature changes, removed fields)
   - Timing-related issues (effects firing too often, infinite loops)
   - Memory leaks or listeners not cleaned up
3. Assume **realistic runtime conditions**, such as:
   - Slow or failed network requests
   - Partial or malformed API responses
   - Rapid user interactions
   - Different browser behaviors
   - Mobile vs desktop constraints
4. For each issue you find, clearly specify:
   - **Location** (file and changed lines)
   - **What runtime error could occur**
   - **Why it would occur**
   - **When it would surface** (user action, lifecycle, edge case)
   - **How severe it is** (crash, silent failure, degraded UX)
   - **Concrete recommendation** to fix or harden the code

### Output format
- Use clear headings
- One issue per section
- Be precise and technical — no generic advice

### Additional Section: Brutal Code Review
Add a final section titled **“Brutal Code Review”** where you:
- Call out risky assumptions
- Highlight fragile logic
- Point out poor defensive coding
- Suggest more robust or idiomatic JavaScript patterns
- Be direct and unapologetically honest

Git diff:
<INSERT_DIFF_BELOW>
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
