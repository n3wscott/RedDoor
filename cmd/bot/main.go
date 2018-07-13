package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/n3wscott/reddoor/pkg/handlers"
	"github.com/n3wscott/reddoor/pkg/listeners"
	"github.com/n3wscott/slack"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"8080"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID"` // required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID"` // required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)

	if env.ChannelID == "" {
		channels, err := client.GetChannels(false)
		if err != nil {
			log.Printf("channel error: %v", err)
		} else {
			log.Printf("channels: %+v", channels)
		}
		return 0
	}

	if env.BotID == "" {
		botid := ""
		bot, err := client.GetBotInfo(&botid)
		if err != nil {
			log.Printf("bot error: %v", err)
		} else {
			log.Printf("bot: %+v", bot)
		}
		return 0
	}

	slackListener := &listeners.SlackListener{
		Client:    client,
		BotID:     env.BotID,
		ChannelID: env.ChannelID,
	}
	go slackListener.ListenAndResponse()

	// Register handler to receive interactive message
	// responses from slack (kicked by user action)
	http.Handle("/interaction", handlers.InteractionHandler{
		VerificationToken: env.VerificationToken,
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
