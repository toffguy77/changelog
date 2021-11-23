package diff

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/toffguy77/changelog/pkg/logger"

	"github.com/toffguy77/changelog/pkg/cvs/repo"

	gogs "github.com/gogs/git-module"
)

func FormatDiff(ctx *context.Context, r *repo.Repository, commits []*gogs.Commit) ([]string, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	var (
		rePullRequest = regexp.MustCompile(`^(Merge pull request #|Pull request #)`)
		reTicket      = regexp.MustCompile(`[\w:/][A-Z]+-[0-9]+`)
	)
	var changes []string
	for _, commit := range commits {
		if rePullRequest.FindString(commit.Message) == "" {
			continue
		}
		ticketID := reTicket.FindString(commit.Message)
		if ticketID == "" {
			ticketID = "NO_TICKET"
		}
		hash := commit.ID
		commitL, err := commitLink(ctx, r, hash.String())
		if err != nil {
			zLog.Errorf("can't get link for commit: %v", err)
			return nil, err
		}
		change := fmt.Sprintf("%s\t%s", strings.Trim(ticketID, "/"), commitL)
		changes = append(changes, change)
	}
	return changes, nil
}

func commitLink(ctx *context.Context, r *repo.Repository, c string) (string, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	data := strings.TrimSpace((*r).URL)
	data = strings.TrimSuffix(data, ".git")
	dataList := strings.Split(data, "/")

	project := dataList[len(dataList)-2]
	repo := dataList[len(dataList)-1]

	return fmt.Sprintf("https://github.com/%s/%s/commit/%s", project, repo, c), nil
}

func PrintDiff(diff []string) {
	if diff == nil {
		fmt.Println("The are no commit messages with Ticket ID")
		return
	}
	for _, change := range diff {
		fmt.Println(change)
	}
}
