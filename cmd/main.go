package main

import (
	"context"
	"log"

	"github.com/toffguy77/changelog/pkg/cvs/repo"
	"github.com/toffguy77/changelog/pkg/diff"
	"github.com/toffguy77/changelog/pkg/flags"
	"github.com/toffguy77/changelog/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Create logger
	ctx := context.Background()
	err := logger.NewLogger(&ctx)
	if err != nil {
		log.Fatalf("can't create new zap logger: %v", err)
	}
	var loggerName logger.LoggerType = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)

	var (
		dataCtx  flags.DataType = "dataType"
		userData                = flags.Data{}
	)
	ctx = context.WithValue(ctx, dataCtx, &userData)
	err = flags.ParseUserFlags(&ctx)
	if err != nil {
		zLog.Fatalf("can't parse user flags: %v", err)
	}
	zLog.Infof("going to work with %s repository starting from the %s start point to the %s", userData.RepoName, userData.StartPoint, userData.EndPoint)

	repository := repo.Repository{
		Name: userData.RepoName,
	}
	err = repository.Clone(&ctx)
	if err != nil {
		zLog.Fatalf("can't clone %s repository: %v", userData.RepoName, err)
		return
	}
	commits, err := repository.Diff(&ctx)
	if err != nil {
		zLog.Fatalf("can't calculate diff: %v", err)
	}
	formattedDiff, err := diff.FormatDiff(&ctx, &repository, commits)
	if err != nil {
		zLog.Fatalf("can't get formatted diff: %v", err)
	}
	diff.PrintDiff(formattedDiff)
}
