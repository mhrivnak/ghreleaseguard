package notify

import (
	"bytes"
	"github.com/mhrivnak/ghreleaseguard/config"
	"log"
	"net/smtp"
	"text/template"
)

func (data *PushData) Send() {
	pushTemplate, err := template.New("push").Parse(pushMessage)
	if err != nil {
		log.Println("notify.PushData.Send: ", err)
		return
	}
	var buf bytes.Buffer
	err = pushTemplate.Execute(&buf, data)
	if err != nil {
		log.Println("notify.PushData.Send: ", err)
		return
	}
	err = smtp.SendMail(config.ServerConfig.SMTPAddress, nil, config.ServerConfig.FromEmail,
		[]string{config.ServerConfig.NotifyEmail}, buf.Bytes())
	if err != nil {
		log.Println("notify.PushData.Send: ", err)
	}
}

func (data *PullRequestData) Send() {
	pushTemplate, err := template.New("pullrequest").Parse(pullRequestMessage)
	if err != nil {
		log.Println("notify.PullRequestData.Send: ", err)
		return
	}
	var buf bytes.Buffer
	err = pushTemplate.Execute(&buf, data)
	if err != nil {
		log.Println("notify.PullRequestData.Send: ", err)
		return
	}
	err = smtp.SendMail(config.ServerConfig.SMTPAddress, nil, config.ServerConfig.FromEmail,
		[]string{config.ServerConfig.NotifyEmail}, buf.Bytes())
	if err != nil {
		log.Println("notify.PullRequestData.Send: ", err)
	}
}

const pushMessage = `
The commit {{.Commit}} was pushed to branch
{{.Branch}}, which appears to be a release branch for version {{.Version}}.
That commit is marked as newer than version {{.Version}} and should probably
not have been merged and pushed.

{{.Url}}

Please investigate ASAP.

Sincerely,
GH Release Guard`

const pullRequestMessage = `
The target branch for the pull request at the below URL appears to be a release
branch for version {{.Version}}. Commit {{.Commit}} in that pull request is
marked as newer than version {{.Version}}, so this pull request may be against
the wrong branch.

{{.Url}}

Please investigate ASAP.

Sincerely,
GH Release Guard`

type PushData struct {
	Branch  string
	Commit  string
	Url     string
	Version string
}

type PullRequestData struct {
	Commit  string
	Url     string
	Version string
}
