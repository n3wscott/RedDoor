package actions

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/n3wscott/slack"
)

// HandleRandomEmoji handles message events.
func (a *Actions) HandleRandomEmoji(ev *slack.MessageEvent) error {


	emoji, err := a.Client.GetEmoji()
	if err != nil {
		return fmt.Errorf("failed to get emojis: %s", err)
	}

	keys := reflect.ValueOf(emoji).MapKeys()
	emote := keys[rand.Intn(len(keys))]


	params := slack.PostMessageParameters{}

	if _, _, err := a.Client.PostMessage(ev.Channel, fmt.Sprintf("%s", emote), params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
