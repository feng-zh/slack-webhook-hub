package hook

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Callback interface {
	OnError(err error)
	OnSuccess(msg string)
}

type Slack struct {
	Tokens []string
}

func (s *Slack) HandleRequest(c *gin.Context, builder Builder) {
	slackToken := c.PostForm("token")
	if s.Tokens != nil && !funk.ContainsString(s.Tokens, slackToken) {
		c.AbortWithStatusJSON(403, "Invalid token")
		log.Println("Invalid token received")
	}
	if responseUrl, ok := c.GetPostForm("response_url"); !ok {
		c.AbortWithStatusJSON(400, "No response_url")
		log.Println("No response_url")
	} else {
		if h, err := builder.NewHooker(c); err != nil {
			c.AbortWithStatusJSON(500, err.Error())
			log.Printf("new hooker failure: %s\n", err)
		} else {
			sr := &SlackRespond{responseUrl: responseUrl}
			if err = h.Hook(sr); err != nil {
				c.AbortWithStatusJSON(500, err.Error())
				log.Printf("trigger hook failure: %s\n", err)
			} else {
				c.JSON(200, "success")
				log.Println("trigger hook success")
			}
		}
	}
}

type SlackRespond struct {
	responseUrl string
}

func (s *SlackRespond) OnError(err error) {
	s.sendRequest(err.Error())
}

func (s *SlackRespond) OnSuccess(msg string) {
	s.sendRequest(msg)
}

func (s *SlackRespond) sendRequest(msg string) {
	json := fmt.Sprintf(`{
	"text": "%s"
}`, msg)
	if req, err := http.NewRequest("POST", s.responseUrl, strings.NewReader(json)); err != nil {
		return
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		defer client.CloseIdleConnections()

		log.Printf("Send response back to slack: %s\n", s.responseUrl)
		if resp, err := client.Do(req); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			defer func() {
				_ = resp.Body.Close()
			}()
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Get slack response: %d, body: %s", resp.StatusCode, string(body))
		}
	}
}
