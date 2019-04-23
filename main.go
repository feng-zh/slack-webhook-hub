package main

import (
	"github.com/feng-zh/slack-webhook-hub/hook"
	"github.com/feng-zh/slack-webhook-hub/hook/travis"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// SLACK_TOKEN
	// TRAVIS_TOKEN
	// GITHUB_REPO
	// BRANCH

	port, ok := os.LookupEnv("PORT")

	if !ok {
		log.Println("$PORT must be set")
		port = "8080"
	}

	hookBuilders := map[string]hook.Builder{"travis": new(travis.Travis)}

	var slackTokens []string
	if slackToken, ok := os.LookupEnv("SLACK_TOKEN"); ok {
		slackTokens = strings.Split(slackToken, ",")
		for i, t := range slackTokens {
			slackTokens[i] = strings.TrimSpace(t)
		}
	}
	s := &hook.Slack{Tokens: slackTokens}

	r := gin.Default()

	r.POST("/:hook", func(c *gin.Context) {
		var h string
		if err := c.ShouldBindUri(&h); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if builder, ok := hookBuilders[h]; ok {
			s.HandleRequest(c, builder)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
			log.Printf("no hook for url: %v\n", c.Request.RequestURI)
		}
	})

	_ = r.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
