package fragbot

import (
	"bytes"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "regexp"
	"strings"
)

func removeRankPrefix(username string) string {
    return regexp.MustCompile(`^\[.*\] `).ReplaceAllString(username, "")
}

func partyInviteGetUsername(msg string) string {
    username := strings.Split(msg, "\n")[1]
    username = strings.Split(username, " has invited you to join their party!")[0]
    return removeRankPrefix(username)
}

func sendToWebhook(webhook string, message string) error {
    if webhook == "" {
        return errors.New("webhook URL is empty")
    }
    payload, err := json.Marshal(map[string]string{"content": message})
    if err != nil {
        return err
    }
    resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func log(args ...interface{}) {
	wh := args[0].(string)
	args = args[1:]
    fmt.Println(args...)
    message := fmt.Sprint(args...)
	if wh != "" {
		err := sendToWebhook(wh, message)
		if err != nil {
			fmt.Println("failed to send webhook:", err)
		}
	}
    
}
