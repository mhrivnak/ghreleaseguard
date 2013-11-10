package notify

import (
	"bytes"
	"github.com/mhrivnak/ghreleaseguard/config"
	"log"
	"net/smtp"
	"text/template"
)

func (data *MessageData) Send(body string) {
	tmpl, err := template.New("message").Parse(body)
	if err != nil {
		log.Println("notify.MessageData.Send: ", err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Println("notify.MessageData.Send: ", err)
		return
	}
	err = smtp.SendMail(config.ServerConfig.SMTPAddress, nil, config.ServerConfig.FromEmail,
		[]string{config.ServerConfig.NotifyEmail}, buf.Bytes())
	if err != nil {
		log.Println("notify.MessageData.Send: ", err)
	}
}

const PushMessage = `
The commit {{.Commit}} was pushed to branch
{{.Branch}}, which appears to be a release branch for version {{.Version}}.
That commit is marked as newer than version {{.Version}} and should probably
not have been merged and pushed.

{{.Url}}

Please investigate ASAP.

Sincerely,
GH Release Guard`

const PullRequestMessage = `
Target branch {{.Branch}} for the pull request at the below URL appears to be a
release branch for version {{.Version}}. Commit {{.Commit}} in that pull
request is marked as newer than version {{.Version}}, so this pull request may
be against the wrong branch.

{{.Url}}

Please investigate ASAP.

Sincerely,
GH Release Guard`

type MessageData struct {
	Branch  string
	Commit  string
	Url     string
	Version string
}
