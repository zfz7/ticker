package apis

import (
	"github.com/gregdel/pushover"
	"os"
)

const PUSHOVER_APP_KEY = "PUSHOVER_APP_KEY"

type NotificationApi interface {
	Send(message *pushover.Message) (pushover.Response, error)
}

type notificationApi struct {
	pushoverClient *pushover.Pushover
	recipientToken string
}

func NewNotificationApi(recipientToken string) NotificationApi {
	pushoverClient := pushover.New(os.Getenv(PUSHOVER_APP_KEY))

	return &notificationApi{
		pushoverClient: pushoverClient,
		recipientToken: recipientToken,
	}
}

func (notificationApi *notificationApi) Send(message *pushover.Message) (pushover.Response, error) {
	recipient := pushover.NewRecipient(notificationApi.recipientToken)
	response, err := notificationApi.pushoverClient.SendMessage(message, recipient)
	if err != nil {
		return pushover.Response{}, err
	}
	return *response, nil
}
