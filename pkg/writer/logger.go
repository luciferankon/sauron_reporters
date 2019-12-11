package writer

import (
	"strings"
)

func (dbWriter MongoDbWriter) logWrite(flowID string) {
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.WriteString("Wrote to DB\n")
	builder.WriteString("-----------\n")
	builder.WriteString("URI: " + dbWriter.URI + "\n")
	builder.WriteString("DB: " + dbWriter.DB + "\n")
	builder.WriteString("Table: " + dbWriter.Table + "\n")
	builder.WriteString("Flow ID: " + flowID + "\n")
	builder.WriteString("---\n")
	dbWriter.Logger.Println(builder.String())
}
