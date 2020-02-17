package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Result(slackWebHook, msg string) {
	if slackWebHook != "" {
		values := map[string]string{"text": fmt.Sprintf("%s\n", msg)}
		jsonValue, _ := json.Marshal(values)
		resp, err := http.Post(slackWebHook, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			Error(err.Error(), nil)
		} else if resp.StatusCode != http.StatusOK {
			Error(fmt.Sprintf("cannot send message to slack [status %d]", resp.StatusCode), nil)
		}
	}
	Info(msg, nil)
}
