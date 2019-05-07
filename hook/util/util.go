package util

import (
	"errors"
	"github.com/feng-zh/slack-webhook-hub/hook"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func GetParameter(c *gin.Context, queryName string, envName string) (string, bool) {
	if value, ok := c.GetQuery(queryName); ok {
		return value, ok
	}
	if value, ok := os.LookupEnv(envName); ok {
		return value, ok
	}
	return "", false
}

func SendRequest(req *http.Request, c hook.Callback) {
	log.Printf("Send request to uri: %s\n", req.URL)

	dryRun, _ := os.LookupEnv("DRY_RUN")
	if dryRun == "true" {
		log.Printf("DEBUG: %s %s\n", req.Method, req.URL)
		for k, v := range req.Header {
			log.Printf("DEBUG: HEADER %s: %s", k, v)
		}
		c.OnSuccess("dry run successful")
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
