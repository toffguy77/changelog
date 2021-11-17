package main

import (
	"context"
	"log"
	"newChangelog/pkg/cvs/repo"
	"newChangelog/pkg/diff"
	"newChangelog/pkg/flags"
	"newChangelog/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Create logger
	ctx := context.Background()
	err := logger.NewLogger(&ctx)
	if err != nil {
		log.Fatalf("can't create new zap logger: %v", err)
	}
	var (
		loggerName logger.LoggerType
	)
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)

	// Populate variables with user input
	var (
		repoName   string
		mode       string
		startPoint string
		endPoint   string
	)
	err = flags.ParseUserFlags(ctx, &repoName, &mode, &startPoint, &endPoint)
	if err != nil {
		zLog.Fatalf("can't parse user flags: %v", err)
	}
	zLog.Infof("going to work with %s repository starting from the %s start point to the %s\n", repoName, startPoint, endPoint)

	repository := repo.Repository{
		Name: repoName,
	}
	err = repository.Clone(ctx, repoName)
	if err != nil {
		zLog.Fatalf("can't clone %s repository: %v", repoName, err)
		return
	}
	commits, err := repository.DiffTags(ctx, startPoint, endPoint)
	if err != nil {
		zLog.Fatalf("can't commits tags: %v", err)
	}
	diff.FormatDiff(ctx, commits)
	//fmt.Println(commits)
}
