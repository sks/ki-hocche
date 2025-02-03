package scmscanner

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/traverse"
	"github.com/sks/kihocche/pkg/logger"
	"github.com/sks/kihocche/pkg/models"
	"github.com/sks/kihocche/pkg/trouble"
	"golang.org/x/sync/errgroup"
)

type githubScanner struct {
	baseScanner
	config SCMConfig
	client *scm.Client
}

func newGithubScanner(config SCMConfig) (Scanner, error) {
	if config.Api == "" {
		config.Api = "https://api.github.com"
	}
	client, err := github.New(config.Api)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	client.Client = config.client()

	return githubScanner{
		config: config,
		client: client,
	}, nil
}

func (s githubScanner) Get(ctx context.Context, filter Filter) (models.Events, error) {
	errGroup, ctx := errgroup.WithContext(ctx)
	events := make(models.Events, 0)
	logger := logger.GetLogger(ctx).With("provider", "github")
	defer func(startTime time.Time) {
		logger.DebugContext(ctx, "scan completed", "duration", time.Since(startTime).String())
	}(time.Now())
	repos, err := traverse.ReposV2(ctx, s.client, scm.ListOptions{
		Page:    1,
		MaxPage: 500,
		Size:    100,
	})
	if err != nil {
		logger.Error("failed to get repos", "error", err)
		return nil, trouble.New("failed to get repos", http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	logger.Debug("repos", "count", len(repos))
	for _, repo := range repos {
		repo := repo
		errGroup.Go(func() error {
			repoEvents, err := s.getRepoEvents(ctx, repo, filter)
			if err != nil {
				return fmt.Errorf("failed to get repo events: %w", err)
			}
			if len(repoEvents) != 0 {
				logger.DebugContext(ctx, "repo events", "repo", repo.Name, "events", len(repoEvents))
			}
			events = append(events, repoEvents...)
			return nil
		})
	}
	err = errGroup.Wait()
	sort.Sort(events)
	return events, err
}

func (s githubScanner) getRepoEvents(ctx context.Context, repo *scm.Repository, filter Filter) (models.Events, error) {
	logger := logger.GetLogger(ctx).With("repo", repo.Name)
	if !filter.isInterestedIn(repo) {
		logger.Debug("skipping repo")
		return nil, nil
	}
	events := make(models.Events, 0)
	for _, event := range filter.Events {
		switch strings.ToLower(event) {
		case string(models.EventTypePush):
			commits, err := s.baseScanner.getCommits(ctx, s.client,
				fmt.Sprintf("%s/%s", repo.Namespace, repo.Name),
				filter, repo)
			if err != nil {
				return nil, fmt.Errorf("failed to get commits: %w", err)
			}
			events = append(events, commits...)
		case string(models.EventTypePullRequest):
			// get pull request events
			pullRequests, err := s.baseScanner.pullRequests(ctx, s.client, fmt.Sprintf("%s/%s", repo.Namespace, repo.Name), filter, repo)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull requests: %w", err)
			}
			events = append(events, pullRequests...)
		case string(models.EventTypeRelease):
			// get release events
			releases, err := s.baseScanner.releases(ctx, s.client, fmt.Sprintf("%s/%s", repo.Namespace, repo.Name), filter, repo)
			if err != nil {
				return nil, fmt.Errorf("failed to get releases: %w", err)
			}
			events = append(events, releases...)
		}
	}
	return events, nil
}
