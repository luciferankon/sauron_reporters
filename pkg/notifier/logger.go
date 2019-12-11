package notifier

import (
	"strings"
)

func (sn SlackNotifier) logNotification(message, recipient string) {
	var builder strings.Builder
	builder.WriteString("\n---\n")
	builder.WriteString("Notified : " + recipient + "\n")
	builder.WriteString("Message : " + message + "\n")
	builder.WriteString("---\n")
	sn.Logger.Println(builder.String())
}

func (sn SlackNotifier) logError(message string, err error) {
	var builder strings.Builder
	builder.WriteString("\n---\n")
	builder.WriteString("Error\n")
	builder.WriteString(message + "\n" + err.Error() + "\n")
	builder.WriteString("---\n")
	sn.Logger.Println(builder.String())
}
