package flags

import (
	"context"
	"errors"
	"flag"
	"strings"

	"github.com/toffguy77/changelog/pkg/logger"
)

type Data struct {
	RepoName   string
	Path       string
	Branch     string
	StartPoint string
	EndPoint   string
}

type DataType string

// Parse user input's flags
func ParseUserFlags(ctx *context.Context) error {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()
	userData := GetData(ctx)

	flag.StringVar(&userData.RepoName, "repo", "", "repository name")
	flag.StringVar(&userData.Path, "path", "", "repository local path")
	flag.StringVar(&userData.Branch, "branch", "", "selected branch to checkout")
	flag.StringVar(&userData.StartPoint, "from", "", "start point for diff")
	flag.StringVar(&userData.EndPoint, "to", "HEAD", "end point for diff")
	//TODO: add debug key
	flag.Parse()

	// Parse repositories from input
	if userData.RepoName == "" && userData.Path == "" {
		err := errors.New("repository name or local path should be provided")
		zLog.Errorf("can't parse user input for repository: %v", err)
		return err
	}

	if userData.RepoName == "" && userData.Path != "" {
		data := strings.Split(userData.Path, "/")
		userData.RepoName = data[len(data)-1]
	}

	// Check from tag input
	if userData.StartPoint == "" {
		err := errors.New("no start point was provided")
		zLog.Errorf("can't parse user input for 'from' start point: %v", err)
		return err
	}

	return nil
}

func GetData(ctx *context.Context) *Data {
	var dataCtx DataType = "dataType"
	userData := (*ctx).Value(dataCtx).(*Data)
	return userData
}
