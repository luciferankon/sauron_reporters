package writer

type Writer interface {
	Write(events map[string]interface{})
}
