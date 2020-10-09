package router

import (
	"github.com/ResultadosDigitais/x9/cmd/worker/git"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
)

type Router interface {
	Healthcheck(c echo.Context) error
	Event(c echo.Context) error
}

type Handler struct {
	Process chan *github.PullRequestEvent
	Session *git.GithubSession
}
