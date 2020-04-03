package git

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type GithubSession struct {
	AccountName  string
	AccountToken string
	Timeout      int
	Context      context.Context
	Client       *github.Client
}

//ParseAuthConfig get the github accounts from environment
func (gs *GithubSession) ParseAuthConfig() error {
	account, ok := os.LookupEnv("GITHUB_ACCOUNTS")
	if !ok {
		return errors.New("No accounts found")
	}

	credentials := strings.Split(account, ":")
	if len(credentials) != 2 {
		return errors.New("Invalid format")
	}
	gs.AccountName = credentials[0]
	gs.AccountToken = credentials[1]
	gs.Context = context.Background()

	return nil
}

func (gs *GithubSession) InitClient() error {
	if err := gs.ParseAuthConfig(); err != nil {
		return err
	}

	staticTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gs.AccountToken})
	client := oauth2.NewClient(gs.Context, staticTokenSource)

	githubClient := github.NewClient(client)
	githubClient.UserAgent = "x9 v1.0"
	_, _, err := githubClient.Users.Get(gs.Context, "")
	if err != nil {
		return err
	}
	gs.Client = githubClient

	gs.Timeout = 60
	if timeout, ok := os.LookupEnv("CLONE_TIMEOUT"); ok {
		if intTimeout, err := strconv.Atoi(timeout); err == nil {
			gs.Timeout = intTimeout
		}
	}
	return nil
}

func (gs *GithubSession) CloneRepository(url, dir string) (*git.Repository, error) {
	localCtx, cancel := context.WithTimeout(gs.Context, time.Duration(gs.Timeout)*time.Second)
	defer cancel()
	auth := &http.BasicAuth{Username: gs.AccountName, Password: gs.AccountToken}

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

func (gs *GithubSession) GetRepository(id int64) (*github.Repository, error) {

	client := gs.Client
	repo, _, err := client.Repositories.GetByID(gs.Context, id)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (gs *GithubSession) OpenIssue(owner, repo, title, message string, labels []string) (string, error) {
	client := gs.Client
	issue := github.IssueRequest{
		Title:  &title,
		Body:   &message,
		Labels: &labels,
	}
	i, _, err := client.Issues.Create(gs.Context, owner, repo, &issue)
	if i == nil {
		return "", err
	}
	return *i.URL, err
}
