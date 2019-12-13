package notifierClient

import (
	"github.com/nlopes/slack"
)

type SlackNotifierClient struct {
	Client *slack.Client
}

func (s SlackNotifierClient) GetUsers() ([]User, error) {
	usrs, err := s.Client.GetUsers()
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0)
	for _, usr := range usrs {
		user := User{
			RealName: usr.RealName,
			Name:     usr.Name,
			ID:       usr.ID,
		}
		users = append(users, user)
	}
	return users, nil
}

func (s SlackNotifierClient) OpenIMChannel(userID string) (bool, bool, string, error) {
	_, _, channelID, err := s.Client.OpenIMChannel(userID)
	if err != nil {
		return false, false, "", err
	}
	return true, true, channelID, nil
}

func (s SlackNotifierClient) SendMessage(channelID string, message string) (string, string, string, error) {
	m := slack.MsgOptionText(message, true)
	_, _, _, err := s.Client.SendMessage(channelID, m)
	if err != nil {
		return "", "", "", err
	}
	return "", "", "", nil
}
