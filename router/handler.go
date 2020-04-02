package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ResultadosDigitais/x9/log"
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
		log.Error("Request error", map[string]interface{}{"error": err.Error()})
		return c.NoContent(http.StatusForbidden)
	}

	event, err := github.ParseWebHook(github.WebHookType(c.Request()), payload)
	if err != nil {
		log.Error("Request error", map[string]interface{}{"error": err.Error()})
		return c.NoContent(http.StatusBadRequest)
	}
	src := c.Request().Header.Get("X-Forward-For")
	switch event := event.(type) {
	case *github.PullRequestEvent:
		if isX9Action(*event.Action) {
			log.Info(fmt.Sprintf("Event received: %s from repository %s", *event.Action, *event.GetRepo().FullName), map[string]interface{}{
				"src_ip": src,
			})
			h.Process <- event
		}

	default:
		log.Warn(fmt.Sprintf("Unexpected event received: %s", event), map[string]interface{}{
			"src_ip": src,
		})
	}
	return c.NoContent(http.StatusOK)

}

func isX9Action(action string) bool {
	listActions := []string{"opened", "edited", "reopened"}
	for _, v := range listActions {
		if action == v {
			return true
		}
	}
	return false
}