package scmscanner

import (
	"context"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/sks/kihocche/pkg/models"
)

type Filter struct {
	Since     time.Duration `yaml:"since"`
	Repos     []string      `yaml:"repos"`
	Namespace []string      `yaml:"namespace"`
	Events    []string      `yaml:"events"`
}

func (f Filter) startDate() time.Time {
	return time.Now().Add(-f.Since)
}

func (f Filter) Filter(event models.Event) bool {
	return event.CreatedOn.After(f.startDate())
}

func (f Filter) isInterestedIn(repo *scm.Repository) bool {
	// check if the repo is in the list of repos and is part of the namespace
	if len(f.Repos) == 0 && len(f.Namespace) == 0 {
		return true
	}
	// check if the repo is in the list of repos
	isInNamespace := false
	if len(f.Namespace) > 0 {
		namespaceMatches := true
		for _, namespace := range f.Namespace {
			if namespace == repo.Namespace {
				isInNamespace = true
				break
			}
		}
		if !namespaceMatches {
			return false // Namespace doesn't match.
		}
	}
	if isInNamespace && len(f.Repos) == 0 {
		return true // Repository is in the namespace and no specific repos are provided.
	}
	for _, repoName := range f.Repos {
		if repoName == repo.Name {
			return true // Repository name matches.
		}
	}
	return false // Repository name doesn't match.
}

type Scanner interface {
	Get(ctx context.Context, filter Filter) (models.Events, error)
}
