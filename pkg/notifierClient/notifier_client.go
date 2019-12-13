package notifierClient

type User struct {
	RealName string
	Name     string
	ID       string
}

type NotifierClient interface {
	GetUsers() ([]User, error)
	OpenIMChannel(user string) (bool, bool, string, error)
	SendMessage(channel string, message string) (string, string, string, error)
}
