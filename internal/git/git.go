package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func GetStagedDiff() (*string, error) {
	cmd := exec.Command("git", "diff", "--staged")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// errMsg := fmt.Sprintf("git diff --staged failed: %v (%s)", err, stderr.String())

		return nil, fmt.Errorf("git diff --staged failed: %v (%s)", err, stderr.String())
	}

	outString := out.String()

	return &outString, nil
}
