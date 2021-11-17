package repo

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/toffguy77/changelog/pkg/logger"

	gogs "github.com/gogs/git-module"
	"go.uber.org/zap"
)

type Repo interface {
	Open() error
	Clone() (string, error)
	Close() error
}

type Repository struct {
	Name string
	Path string
	URL  string
	Data *gogs.Repository
}

func getRepoPath(str string) (string, error) {
	//return "ssh://git@stash.sigma.sbrf.ru:7999/" + str + ".git"
	if strings.HasPrefix(str, "/") {
		str = str[1:]
	}
	if strings.HasSuffix(str, ".git") {
		if len(str) < 4 {
			return "", errors.New("repository name is invalid")
		}
		str = str[:len(str)-4]
	}
	return "git@github.com:toffguy77/" + str + ".git", nil
}

func (r *Repository) Clone(ctx context.Context, path string) error {
	var loggerName logger.LoggerType
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)
	defer zLog.Sync()

	pathFrom, err := getRepoPath(path)
	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir("", "changelog")
	if err != nil {
		zLog.Errorf("can't obtain temporary directory to check out repository: %v", err)
		return err
	}
	zLog.Infof("clone %s to %s directory", pathFrom, tempDir)

	pathTo := filepath.Join(tempDir, path)

	if _, err := os.Stat(pathTo); !os.IsNotExist(err) {
		zLog.Infof("removing temporary directory %s", pathTo)
		err := os.RemoveAll(pathTo)
		if err != nil {
			zLog.Errorf("can't delete temporary directory: %v", err)
			return err
		}
		return fmt.Errorf("temporary directory at %s was not created", pathTo)
	}

	options := gogs.CloneOptions{
		Depth:   100,
		Timeout: 300 * time.Second,
	}

	err = gogs.Clone(pathFrom, pathTo, options)
	if err != nil {
		zLog.Errorf("can't clone %s repository: %v", pathFrom, err)
		return err
	}

	cmd := gogs.NewCommand("fetch")
	cmd.AddArgs("--tags")
	_, err = cmd.RunInDir(pathTo)
	if err != nil {
		zLog.Error("can't fetch tags")
		return err
	}

	r.Path = pathTo
	r.URL = pathFrom
	r.Data, err = gogs.Open(r.Path)
	return err
}

func (r *Repository) DiffTags(ctx context.Context, tagStart, tagEnd string) ([]*gogs.Commit, error) {
	var loggerName logger.LoggerType
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)
	defer zLog.Sync()

	commitStart, err := r.GetCommit(ctx, tagStart)
	if err != nil {
		return nil, err
	}
	commitEnd, err := r.GetCommit(ctx, tagEnd)
	if err != nil {
		return nil, err
	}
	commitsDiff, err := r.DiffCommits(ctx, commitStart, commitEnd)
	if err != nil {
		zLog.Errorf("Can't get diff between commits %v and %v", commitStart, commitEnd)
		return nil, err
	}

	return commitsDiff, nil
}

func (r *Repository) DiffCommits(ctx context.Context, commitStart, commitEnd string) ([]*gogs.Commit, error) {
	var loggerName logger.LoggerType
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)
	defer zLog.Sync()

	options := gogs.LogOptions{
		Timeout: 300 * time.Second,
	}
	rev := commitStart + ".." + commitEnd
	diff, err := r.Data.Log(rev, options)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func (r *Repository) GetCommit(ctx context.Context, tag string) (string, error) {
	var loggerName logger.LoggerType
	loggerName = "ZapLogger"
	zLog := ctx.Value(loggerName).(*zap.SugaredLogger)
	defer zLog.Sync()

	if tag == "latest" {
		return "HEAD", nil
	}

	id, err := r.Data.TagCommitID(tag)
	if err != nil {
		zLog.Errorf("can't find commit id for tag %s", tag)
		return "", err
	}
	return id, nil
}
