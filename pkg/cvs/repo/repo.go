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

	"github.com/toffguy77/changelog/pkg/flags"
	"github.com/toffguy77/changelog/pkg/logger"

	gogs "github.com/gogs/git-module"
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

const DEFAUTL_TIMEOUT = 600 * time.Second

func getRepoPath(str string) (string, error) {
	str = strings.TrimPrefix(str, "/")
	if strings.HasSuffix(str, ".git") {
		if len(str) < 4 {
			return "", errors.New("repository name is invalid")
		}
		str = str[:len(str)-4]
	}
	return "git@github.com:" + str + ".git", nil
}

func (r *Repository) Clone(ctx *context.Context) error {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	userData := flags.GetData(ctx)

	pathFrom, err := getRepoPath(userData.RepoName)
	if err != nil {
		zLog.Errorf("error getting repository path: %v", err)
		return err
	}

	if userData.Path == "" {
		tempDir, err := ioutil.TempDir("", "changelog")
		if err != nil {
			zLog.Errorf("can't obtain temporary directory to check out repository: %v", err)
			return err
		}
		zLog.Infof("clone %s to %s directory", pathFrom, tempDir)

		path := strings.SplitN(userData.RepoName, "/", 2)[1]
		pathTo := filepath.Join(tempDir, path)
		r.Path = pathTo

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
			Timeout: DEFAUTL_TIMEOUT,
		}
		if userData.Branch != "" {
			options.Branch = userData.Branch
		}

		err = gogs.Clone(pathFrom, pathTo, options)
		if err != nil {
			zLog.Errorf("can't clone %s repository: %v", pathFrom, err)
			return err
		}
		userData.Path = pathTo
	} else {
		r.Path = userData.Path
	}

	r.Data, err = gogs.Open(r.Path)
	if err != nil {
		zLog.Errorf("can't open repository %s", r.Path)
		return err
	}

	repoURL, err := GetRepoURL(ctx)
	if err != nil {
		zLog.Errorf("can't find repo url: %v", err)
		return nil
	}
	r.URL = repoURL

	cmdBuff := gogs.NewCommand("config")
	cmdBuff.AddArgs("http.postBuffer", "524288000")
	zLog.Infof("running git comand: %s", cmdBuff.String())
	_, err = cmdBuff.RunInDir(r.Path)
	if err != nil {
		zLog.Errorf("can't set postBuffer: %v", err)
		return err
	}

	zLog.Infof("git pull @ %v", r.Path)
	r.Data.Pull(gogs.PullOptions{Timeout: DEFAUTL_TIMEOUT})

	cmd := gogs.NewCommand("fetch")
	cmd.AddArgs("--tags")
	cmd.AddArgs("-f")
	zLog.Infof("running git comand: %s", cmd.String())
	_, err = cmd.RunInDirWithTimeout(DEFAUTL_TIMEOUT, r.Path)
	if err != nil {
		zLog.Errorf("can't fetch tags: %v", err)
		return err
	}

	zLog.Info("repository clone completed")
	return nil
}

func (r *Repository) Diff(ctx *context.Context) ([]*gogs.Commit, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	userData := flags.GetData(ctx)

	if tag, err := r.GetCommit(ctx, userData.StartPoint); tag != "" && err == nil {
		commits, err := r.DiffTags(ctx)
		if err != nil {
			zLog.Errorf("can't diff tags: %v", err)
			return commits, err
		}
		zLog.Infof("calculating diff for %s..%s tags", userData.StartPoint, userData.EndPoint)
		return commits, err
	}
	commits, err := r.DiffCommits(ctx, userData.StartPoint, userData.EndPoint)
	if err != nil {
		zLog.Errorf("can't diff commits: %v", err)
		return commits, err
	}
	zLog.Infof("calculating diff for %s..%s commits", userData.StartPoint, userData.EndPoint)
	return commits, err
}

func (r *Repository) DiffTags(ctx *context.Context) ([]*gogs.Commit, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	userData := flags.GetData(ctx)

	commitStart, err := r.GetCommit(ctx, userData.StartPoint)
	if err != nil {
		zLog.Error("can't find tag %s", userData.StartPoint)
		return nil, err
	}
	commitEnd, err := r.GetCommit(ctx, userData.EndPoint)
	if err != nil {
		zLog.Error("can't find tag %s", userData.EndPoint)
		return nil, err
	}
	commitsDiff, err := r.DiffCommits(ctx, commitStart, commitEnd)
	if err != nil {
		zLog.Errorf("Can't get diff between commits %v and %v", commitStart, commitEnd)
		return nil, err
	}

	return commitsDiff, nil
}

func (r *Repository) DiffCommits(ctx *context.Context, commitStart, commitEnd string) ([]*gogs.Commit, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	options := gogs.LogOptions{
		Timeout: DEFAUTL_TIMEOUT,
	}
	rev := commitStart + ".." + commitEnd
	diff, err := r.Data.Log(rev, options)
	if err != nil {
		zLog.Errorf("can't perform git log for %s: %v", rev, err)
		return nil, err
	}

	return diff, nil
}

func (r *Repository) GetCommit(ctx *context.Context, tag string) (string, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	if tag == "HEAD" {
		return tag, nil
	}

	id, err := r.Data.TagCommitID(tag)
	if err != nil {
		zLog.Errorf("can't find commit id for tag %s", tag)
		return "", err
	}
	return id, nil
}

func GetRepoURL(ctx *context.Context) (string, error) {
	zLog := logger.GetLogger(ctx)
	defer zLog.Sync()

	userData := flags.GetData(ctx)

	cmd := gogs.NewCommand("config")
	cmd.AddArgs("--get")
	cmd.AddArgs("remote.origin.url")
	zLog.Infof("running git comand: %s", cmd.String())
	zLog.Infof("directory to execute command: %s", userData.Path)
	dataByte, err := cmd.RunInDir(userData.Path)
	if err != nil {
		zLog.Errorf("can't get repo url: %v", err)
		return "", err
	}
	data := strings.TrimSpace(string(dataByte))

	return data, nil
}
