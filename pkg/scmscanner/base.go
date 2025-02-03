package scmscanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/sks/kihocche/pkg/models"
)

type baseScanner struct {
}

func (baseScanner) commitToEvents(repo *scm.Repository, commits []*scm.Commit) models.Events {
	events := make(models.Events, len(commits))
	for i, commit := range commits {
		msgs := strings.Split(commit.Message, "\n")
		events[i] = models.Event{
			ID:          commit.Sha,
			Name:        msgs[0],
			Link:        commit.Link,
			Description: strings.Join(msgs[0:], "\n"),
			Type:        models.EventTypePush,
			Actor: models.EventActor{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			},
			Repo: models.Repo{
				Name: repo.Name,
				URL:  repo.Link,
			},
			CreatedOn: commit.Author.Date,
		}
	}
	return events
}

func (b baseScanner) getCommits(
	ctx context.Context,
	client *scm.Client,
	repoID string,
	filter Filter,
	repo *scm.Repository,
) (models.Events, error) {
	events := make(models.Events, 0)
	page := 0
	for {
		page++
		commits, _, err := client.Git.ListCommits(ctx, repoID, scm.CommitListOptions{
			Page: page,
			Size: 100,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}
		if len(commits) == 0 {
			break
		}
		events = events.Add(b.commitToEvents(repo, commits), filter.Filter)
		if commits[len(commits)-1].Author.Date.Before(filter.startDate()) {
			break
		}
	}
	return events, nil
}

func (b baseScanner) pullRequests(
	ctx context.Context,
	client *scm.Client,
	repoID string,
	filter Filter,
	repo *scm.Repository,
) (models.Events, error) {
	events := make(models.Events, 0)
	page := 0
	for {
		page++
		pullRequests, _, err := client.PullRequests.List(ctx, repoID, scm.PullRequestListOptions{
			Page: page,
			Size: 100,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}
		if len(pullRequests) == 0 {
			break
		}
		events = events.Add(b.pullToEvents(repo, pullRequests), filter.Filter)
		if pullRequests[len(pullRequests)-1].Updated.Before(filter.startDate()) {
			break
		}
	}
	return events, nil
}

func (baseScanner) pullToEvents(repo *scm.Repository, pullRequests []*scm.PullRequest) models.Events {
	events := make(models.Events, len(pullRequests))
	for i, pullRequest := range pullRequests {
		events[i] = models.Event{
			ID:          pullRequest.Sha,
			Name:        pullRequest.Title,
			Link:        pullRequest.Link,
			Description: pullRequest.Body,
			Type:        models.EventTypePullRequest,
			Actor: models.EventActor{
				Name:  pullRequest.Author.Name,
				Email: pullRequest.Author.Email,
			},
			Repo: models.Repo{
				Name: repo.Name,
				URL:  repo.Link,
			},
			CreatedOn: pullRequest.Created,
		}
	}
	return events
}

func (b baseScanner) releases(
	ctx context.Context,
	client *scm.Client,
	repoID string,
	filter Filter,
	repo *scm.Repository,
) (models.Events, error) {
	events := make(models.Events, 0)
	page := 0
	for {
		page++
		pullRequests, _, err := client.Releases.List(ctx, repoID, scm.ReleaseListOptions{
			Page: page,
			Size: 100,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}
		if len(pullRequests) == 0 {
			break
		}
		events = events.Add(b.releasesToEvents(repo, pullRequests), filter.Filter)
		if pullRequests[len(pullRequests)-1].Published.Before(filter.startDate()) {
			break
		}
	}
	return events, nil
}

func (baseScanner) releasesToEvents(repo *scm.Repository, releases []*scm.Release) models.Events {
	events := make(models.Events, len(releases))
	for i, release := range releases {
		events[i] = models.Event{
			ID:          fmt.Sprintf("%d", release.ID),
			Name:        release.Title,
			Link:        release.Link,
			Description: release.Description,
			Type:        models.EventTypeRelease,
			Repo: models.Repo{
				Name: repo.Name,
				URL:  repo.Link,
			},
			CreatedOn: release.Published,
		}
	}
	return events
}
