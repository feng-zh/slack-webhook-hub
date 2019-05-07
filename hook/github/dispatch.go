package github

import (
	"fmt"
	"github.com/feng-zh/slack-webhook-hub/hook"
	. "github.com/feng-zh/slack-webhook-hub/hook/util"
	"log"
	"net/http"
	"strings"
)

func (t *hooker) dispatches(c hook.Callback) error {
	log.Println("start trigger github RepositoryDispatch event")
	json := `{"event_type": "dispatch"}`
	if req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/dispatches", t.repo), strings.NewReader(json)); err != nil {
		return err
	} else {
		req.SetBasicAuth("", t.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/vnd.github.everest-preview+json")

		go SendRequest(req, c)
		return nil
	}
}
