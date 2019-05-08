package travis

import (
	"errors"
	"fmt"
	"github.com/feng-zh/slack-webhook-hub/hook"
	. "github.com/feng-zh/slack-webhook-hub/hook/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Travis struct {
}

var _ hook.Builder = (*Travis)(nil)

func (t *Travis) NewHooker(c *gin.Context) (hook.Hooker, error) {
	travisToken, ok := GetParameter(c, "travisToken", "TRAVIS_TOKEN")
	if !ok {
		return nil, errors.New("no travis token defined")
	}
	githubRepo, ok := GetParameter(c, "githubRepo", "GITHUB_REPO")
	if !ok {
		return nil, errors.New("no github repo defined")
	}
	branch, ok := GetParameter(c, "branch", "BRANCH")
	if !ok {
		return nil, errors.New("no branch defined")
	}
	h := &hooker{token: travisToken, repo: githubRepo, branch: branch}
	return h, nil
}

type hooker struct {
	token  string
	repo   string
	branch string
}

var _ hook.Hooker = (*hooker)(nil)

func (t *hooker) Hook(c hook.Callback) error {
	log.Println("start trigger travis ci")
	json := fmt.Sprintf(`{
	"request": {
		"message": "trigger from slack",
		"branch": "%s"
	}
}`, t.branch)
	if req, err := http.NewRequest("POST", fmt.Sprintf("https://api.travis-ci.org/repo/%s/requests", strings.ReplaceAll(t.repo, "/", "%2F")), strings.NewReader(json)); err != nil {
		return err
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", t.token))
		req.Header.Set("Travis-API-Version", "3")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		go SendRequest(req, c)
		return nil
	}
}
