package core

import (
	"context"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func CloneRepository(session *Session, url string, dir string) (*git.Repository, error) {
	localCtx, cancel := context.WithTimeout(session.Context, time.Duration(*session.Options.CloneRepositoryTimeout)*time.Second)
	defer cancel()
	auth := &http.BasicAuth{Username: session.Config.GitHubAccessTokens[0].AccountName, Password: session.Config.GitHubAccessTokens[0].Token}

	repository, err := git.PlainCloneContext(localCtx, dir, false, &git.CloneOptions{
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
		URL:               url,
		SingleBranch:      true,
		Tags:              git.NoTags,
		Auth:              auth,
	})

	if err != nil {
		return nil, err
	}

	return repository, nil
}
