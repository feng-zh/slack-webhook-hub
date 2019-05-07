package github

import (
	"errors"
	"github.com/feng-zh/slack-webhook-hub/hook"
	. "github.com/feng-zh/slack-webhook-hub/hook/util"
	"github.com/gin-gonic/gin"
)

type Github struct {
}

var _ hook.Builder = (*Github)(nil)

func (t *Github) NewHooker(c *gin.Context) (hook.Hooker, error) {
	token, ok := GetParameter(c, "token", "GITHUB_TOKEN")
	if !ok {
		return nil, errors.New("no github token defined")
	}
	githubRepo, ok := GetParameter(c, "githubRepo", "GITHUB_REPO")
	if !ok {
		return nil, errors.New("no github repo defined")
	}
	branch, ok := GetParameter(c, "branch", "BRANCH")
	if !ok {
		return nil, errors.New("no branch defined")
	}
	event, ok := GetParameter(c, "event", "EVENT")
	if !ok {
		event = "dispatches"
	}
	h := &hooker{token: token, repo: githubRepo, branch: branch, event: event}
	return h, nil
}

type hooker struct {
	token  string
	repo   string
	branch string
	event  string
}

var _ hook.Hooker = (*hooker)(nil)

func (t *hooker) Hook(c hook.Callback) error {
	return t.dispatches(c)
}
