package examples

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

const (
	appID    = 6
	appHash  = "YOUR_APP_HASH"
	botToken = "YOUR_BOT_TOKEN"
)

func main() {
	// Create a new client
	client, _ := telegram.NewClient(telegram.ClientConfig{
		AppID:    appID,
		AppHash:  appHash,
		LogLevel: telegram.LogInfo,
	})

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Authenticate the client using the bot token
	if err := client.LoginBot(botToken); err != nil {
		panic(err)
	}

	client.On("message", convEventHandler)

	// new conversation
	conv, _ := client.NewConversation("username or id", false, 30) // 30 is the timeout in seconds, false means it's not a private conversation
	defer conv.Close()
	_, err := conv.Respond("Hello, Please reply to this message")
	if err != nil {
		panic(err)
	}
	resp, err := conv.GetResponse() // wait for the response
	// resp, err := conv.GetReply() // wait for the reply
	// conv.MarkRead() // mark the conversation as read
	// conv.WaitEvent() // wait for any custom update
	if err != nil {
		panic(err)
	}

	fmt.Println("Response:", resp.Text())
}

func convEventHandler(m *telegram.NewMessage) error {
	response, err := m.Ask("What's your name?")
	if err != nil {
		return err
	}

	response.Reply("Nice to meet you, " + response.Text())
	return nil
}
