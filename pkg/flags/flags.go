package flags

import (
	"context"
	"errors"
	"flag"
	"newChangelog/pkg/logger"

	"go.uber.org/zap"
)

// Parse user input's flags
func ParseUserFlags(ctx context.Context, repo *string, mode *string, from *string, to *string) error {
	var loggerName logger.LoggerType
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)
	defer zLog.Sync()

	flag.StringVar(repo, "repo", "", "repository name")
	flag.StringVar(mode, "mode", "tags", "calculate diff for tags of commits, tags is default")
	flag.StringVar(from, "from", "", "start point for diff")
	flag.StringVar(to, "to", "latest", "end point for diff")
	//TODO: add debug key
	flag.Parse()

	// Parse repositories from input
	if *repo == "" {
		err := errors.New("no repository name was provided")
		zLog.Errorf("can't parse user input for repository: %v", err)
		return err
	}

	// Check from tag input
	if *from == "" {
		err := errors.New("no start point was provided")
		zLog.Errorf("can't parse user input for 'from' start point: %v", err)
		return err
	}
	return nil
}
