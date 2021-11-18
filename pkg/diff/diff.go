package diff

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	gogs "github.com/gogs/git-module"
)

func FormatDiff(ctx context.Context, commits []*gogs.Commit) []string {
	var re = regexp.MustCompile(`^\w*[A-Z]+-[0-9]+`)
	var changes []string
	for _, commit := range commits {
		ticketID := re.FindString(commit.Message)
		if ticketID == "" {
			ticketID = "NO_TICKET"
		}
		hash := commit.ID
		message := strings.TrimSpace(commit.Message)
		change := fmt.Sprintf("%s\t%s\t%s", hash, ticketID, message)
		changes = append(changes, change)
	}
	return changes
}
