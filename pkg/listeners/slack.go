package listeners

import (
	"fmt"
	"log"
	"strings"

	"github.com/n3wscott/slack"

	"github.com/n3wscott/reddoor/pkg/actions"
)

type SlackListener struct {
	Client    *slack.Client
	BotID     string
	ChannelID string

	action *actions.Actions
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {

	s.action = &actions.Actions{
		Client:    s.Client,
		BotID:     s.BotID,
		ChannelID: s.ChannelID,
	}

	rtm := s.Client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.ChannelID {
		log.Printf("channel mismatch, not %s", s.ChannelID)
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}
	log.Printf("channel match")

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.BotID)) {
		log.Printf("Not prefix %s", fmt.Sprintf("<@%s> ", s.BotID))
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}
	log.Printf("prefix match")

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 {
		return fmt.Errorf("invalid message")
	}

	switch m[0] {
	case "beer":
		return s.action.HandleBeer(ev)
	default:
		return s.action.HandleRandomEmoji(ev)
	}

	return nil
}
