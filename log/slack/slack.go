package slack

import (
	"encoding/json"
	"os"
)

func Send(repository, vulnerability, filename, values, id string) {
	slackWebHook := os.Getenv("SLACK_WEBHOOK")
	if slackWebHook != "" {
		values := map[string]interface{}{
			"text": ":warning: *Vulnerability found*",
			"attachments": []map[string]interface{}{
				map[string]interface{}{
					"color": "#AB3117",
					"fields": []map[string]interface{}{
						map[string]interface{}{
							"title": "Repository",
							"value": repository,
							"short": false,
						},
						map[string]interface{}{
							"title": "Vulnerability",
							"value": vulnerability,
							"short": false,
						},
						map[string]interface{}{
							"title": "File",
							"value": filename,
							"short": false,
						},
						map[string]interface{}{
							"title": "Values",
							"value": values,
							"short": false,
						},
					},
				},
				map[string]interface{}{
					"callback_id":     id,
					"title":           "Actions",
					"color":           "#3AA3E3",
					"attachment_type": "default",
					"actions": []map[string]string{
						map[string]string{
							"name":  "Open Issue",
							"text":  "Open Issue",
							"type":  "button",
							"value": "open_issue",
						},
					},
				},
			},
		}
		json.Marshal(values)
		// resp, err := http.Post(slackWebHook, "application/json", bytes.NewBuffer(jsonValue))
		// if err != nil {
		// 	log.Error(err.Error(), nil)
		// } else if resp.StatusCode != http.StatusOK {
		// 	log.Error(fmt.Sprintf("cannot send message to slack [status %d]", resp.StatusCode), nil)
		// }
	}
}
