package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ResultadosDigitais/x9/cmd/worker/actions"
	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/crypto"
	"github.com/ResultadosDigitais/x9/log"

	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"
)

// HealthCheck returns a status code 200
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "WORKING!")
}

// Event receives an event request, check its authencity and process it
func (h *Handler) Event(c echo.Context) error {
	payload, err := github.ValidatePayload(c.Request(), []byte(config.Opts.GithubSecretWebhook))
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

func (h *Handler) Action(c echo.Context) error {
	secret := config.Opts.SlackSecretKey
	version := "v0"
	timestamp := c.Request().Header.Get("X-Slack-Request-Timestamp")
	slackSignature := c.Request().Header.Get("X-Slack-Signature")
	body, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		log.Error("Error parsing body", map[string]interface{}{"error": err.Error()})
		return c.NoContent(http.StatusBadRequest)
	}

	content := fmt.Sprintf("%s:%s:%s", version, timestamp, string(body))
	if !crypto.ValidateHMAC(content, slackSignature, version, secret) {
		log.Error("Error validating signature", nil)

		return c.NoContent(http.StatusBadRequest)
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Error("Error on parsing form values", map[string]interface{}{"error": err.Error()})
	}
	payload := values.Get("payload")
	actions.ProcessAction(payload, h.Session)
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
