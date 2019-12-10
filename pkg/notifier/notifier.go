package notifier

type Notifier interface {
	Notify(events map[string]interface{})
}