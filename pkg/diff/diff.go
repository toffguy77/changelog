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

func FormatDiff(ctx *context.Context, r *repo.Repository, commits []*gogs.Commit) ([][]string, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	var (
		rePullRequest = regexp.MustCompile(`^(Merge pull request #|Pull request #)`)
		reTicket      = regexp.MustCompile(`[\w:/][A-Z]+-[0-9]+`)
	)
	var changes [][]string
	for _, commit := range commits {
		change := make([]string, 0, 2)
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
		change = append(change, strings.Trim(ticketID, "/"), commitL)
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
	repository := dataList[len(dataList)-1]

	return fmt.Sprintf("https://github.com/%s/%s/commit/%s", project, repository, c), nil
}
