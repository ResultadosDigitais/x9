package actions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ResultadosDigitais/x9/git"
	"github.com/ResultadosDigitais/x9/log"
	"github.com/ResultadosDigitais/x9/management"
)

type SlackAction struct {
	Type            string          `json:"type"`
	Actions         []Action        `json:"actions"`
	CallbackID      string          `json:"callback_id"`
	Team            Team            `json:"team"`
	Channel         Channel         `json:"channel"`
	User            User            `json:"user"`
	ActionTS        string          `json:"action_ts"`
	MessageTS       string          `json:"message_ts"`
	AttachmentID    string          `json:"attachment_id"`
	Token           string          `json:"token"`
	OriginalMessage json.RawMessage `json:"original_message"`
	ResponseURL     string          `json:"response_url"`
	TriggerID       string          `json:"trigger_id"`
}

type Action struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ProcessAction(body string, gs *git.GithubSession) error {
	var action SlackAction
	if err := json.Unmarshal([]byte(body), &action); err != nil {
		log.Error("Error parsing json", map[string]interface{}{"error": err.Error()})
	}

	if action.Type == "interactive_message" &&
		action.Actions[0].Value == "open_issue" {
		vuln, err := management.GetVulnerabilityByID(action.CallbackID)
		if err != nil {
			log.Error("Error on getting vulnerability", map[string]interface{}{"error": err.Error()})
			return err
		}
		if vuln.IssueURL != "" {
			return nil
		}
		title := getIssueTitle(vuln.Name, vuln.FileName)
		body := getIssueBody(vuln.Name, vuln.Value, vuln.FileName)
		labels := []string{
			"X9",
			"Security",
		}
		owner, repo := getRepoInfo(vuln.Repository)
		issueURL, err := gs.OpenIssue(owner, repo, title, body, labels)
		if err != nil {
			log.Error("Error openning issue", map[string]interface{}{"error": err.Error()})
		} else {
			if err := management.SetIssueURL(vuln.ID, issueURL); err != nil {
				log.Error("Error on updating issue URL", map[string]interface{}{"error": err.Error()})
			}
		}

		return err
	}
	return nil
}

func getRepoInfo(url string) (string, string) {
	r := regexp.MustCompile(`((https://([a-z]+)\.com/)|(\.git$))`)
	ownerAndRepo := r.ReplaceAllString(url, "")
	info := strings.Split(ownerAndRepo, "/")
	return info[0], info[1]
}

func getIssueTitle(name, filename string) string {
	return fmt.Sprintf(
		"[Vulnerability] Sensitive data: %s in %s",
		name,
		filename,
	)
}

func getIssueBody(name, values, filename string) string {
	return fmt.Sprintf(
		`# X9 Vulnerability Report
### Description
- **Vulnerability**: %s
- **Values**: %s
- **Filename**: %s`,
		name,
		values,
		filename,
	)
}
