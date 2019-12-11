package eventlistener

import (
	"fmt"
	"strings"
)

func (l Listener) logStart(sName string) {
	var builder strings.Builder
	builder.WriteString("Starting event listener...\n")
	builder.WriteString("---\n")
	builder.WriteString("Reading from stream: " + sName + "\n")
	builder.WriteString("---\n")
	l.Logger.Println(builder.String())
}

func (l Listener) logRead(eventID string) {
	message := fmt.Sprintf("Read event with id: %s", eventID)
	l.Logger.Println(message)
}
