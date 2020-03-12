package router

import (
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/labstack/echo"
)

type Handler struct {
	Process chan *github.PullRequestEvent
}

// HealthCheck returns a status code 200
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!")
}

// Event receives an event request, check its authencity and process it
func (h *Handler) Event(c echo.Context) error {
	payload, err := github.ValidatePayload(c.Request(), []byte(os.Getenv("GITHUB_SECRET_WEBHOOK")))
	if err != nil {
		return c.NoContent(http.StatusForbidden)
	}

	event, err := github.ParseWebHook(github.WebHookType(c.Request()), payload)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	switch event := event.(type) {
	case *github.PullRequestEvent:
		h.Process <- event
	}
	return c.NoContent(http.StatusOK)

}
