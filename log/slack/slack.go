package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ResultadosDigitais/x9/log"
)

func Send(slackWebHook string, fields map[string]interface{}) {
	if slackWebHook != "" {
		slackMessage := formatMessage(fields)
		values := map[string]string{"text": slackMessage}
		jsonValue, _ := json.Marshal(values)
		resp, err := http.Post(slackWebHook, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Error(err.Error(), nil)
		} else if resp.StatusCode != http.StatusOK {
			log.Error(fmt.Sprintf("cannot send message to slack [status %d]", resp.StatusCode), nil)
		}
	}
}

func formatMessage(fields map[string]interface{}) string {

	if _, ok := fields["matches"]; ok {
		return fmt.Sprintf(":warning: *Ooops I found something...*\n"+
			"*Repository:* %s\n*File:* %s\n*Vulnerability:* %s\n*Matches:* %d\n*Values:* %s\n",
			fields["repo"], fields["file"], fields["vuln"], fields["matches"], fields["values"])

	}
	return fmt.Sprintf(":warning: *Ooops I found something...*\n"+
		"*Repository:* %s\n*File:* %s\n*Vulnerability:* %s\n",
		fields["repo"], fields["file"], fields["vuln"])

}
