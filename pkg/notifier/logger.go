package notifier

import (
	"strings"
)

func (sn SlackNotifier) logNotification(message, recipient string) {
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.WriteString("Notified : " + recipient + "\n")
	builder.WriteString("Message : " + message + "\n")
	builder.WriteString("---\n")
	sn.Logger.Println(builder.String())
}
