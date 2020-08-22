package controller

import (
	"dj-api/app/config"
	"dj-api/tools/slack"
	"github.com/gin-gonic/gin"
)

type SentryForm struct {
	ProjectName string `json:"project_name"`
	Message     string `json:"message"`
	Culprit     string `json:"culprit"`
	Url         string `json:"url"`
	Level       string `json:"level"`
}

func Notify(c *gin.Context) {
	var form SentryForm

	if err := c.ShouldBind(&form); err != nil {
		c.Writer.WriteString("ERR")
		return
	}

	attachment1 := slack.Attachment{}
	attachment1.AddField(slack.Field{Title: "Project", Value: form.ProjectName}).
		AddField(slack.Field{Title: "level", Value: form.Level}).
		AddField(slack.Field{Title: "culprit", Value: form.Culprit})
	attachment1.AddAction(slack.Action{Type: "button", Text: "打开", Url: form.Url, Style: "primary"})
	payload := slack.Payload{
		Text:        form.Message,
		Username:    "robot",
		Channel:     "#" + c.DefaultQuery("source", "dj"),
		IconEmoji:   ":cupid:",
		Attachments: []slack.Attachment{attachment1},
	}
	slack.Send(config.C.Slack.HookUrl, "", payload)
	c.Writer.WriteString("OK")
	return
}
