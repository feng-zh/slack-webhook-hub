package travis

import (
	"errors"
	"fmt"
	"github.com/feng-zh/slack-webhook-hub/hook"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

type Travis struct {
}

var _ hook.Builder = (*Travis)(nil)

func (t *Travis) NewHooker(c *gin.Context) (hook.Hooker, error) {
	travisToken, ok := getParameter(c, "travisToken", "TRAVIS_TOKEN")
	if !ok {
		return nil, errors.New("no travis toke defined")
	}
	githubRepo, ok := getParameter(c, "githubRepo", "GITHUB_REPO")
	if !ok {
		return nil, errors.New("no github repo defined")
	}
	branch, ok := getParameter(c, "branch", "BRANCH")
	if !ok {
		return nil, errors.New("no branch defined")
	}
	h := &hooker{token: travisToken, repo: githubRepo, branch: branch}
	return h, nil
}

func getParameter(c *gin.Context, queryName string, envName string) (string, bool) {
	if value, ok := c.GetQuery(queryName); ok {
		return value, ok
	}
	if value, ok := os.LookupEnv(envName); ok {
		return value, ok
	}
	return "", false
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
	if req, err := http.NewRequest("POST", fmt.Sprintf("https://api.travis-ci.org/repo/%s/requests", t.repo), strings.NewReader(json)); err != nil {
		return err
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", t.token))
		req.Header.Set("Travis-API-Version", "3")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		go sendRequest(req, c)
		return nil
	}
}

func sendRequest(req *http.Request, c hook.Callback) {
	log.Printf("Send request to uri: %s\n", req.RequestURI)

	dryRun, _ := os.LookupEnv("DRY_RUN")
	if dryRun == "true" {
		log.Printf("DEBUG: %s %s\n", req.Method, req.RequestURI)
		for k, v := range req.Header {
			log.Printf("DEBUG: HEADER %s: %s", k, v)
		}
		c.OnSuccess("dry run success")
		return
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()
	if resp, err := client.Do(req); err != nil {
		c.OnError(err)
	} else {
		switch {
		case resp.StatusCode < 300:
			c.OnSuccess("command is started")
		default:
			c.OnError(errors.New("command cannot execute"))
		}
	}
}
