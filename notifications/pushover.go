package notifications

import (
	"github.com/aetaric/checkrr/logging"
	"github.com/gregdel/pushover"
	"github.com/spf13/viper"
)

type Pushover struct {
	config        viper.Viper
	AllowedNotifs []string
	apiToken      string
	recipient     *pushover.Recipient
	bot           *pushover.Pushover
	Log           *logging.Log
}

func (p Pushover) Notify(title string, description string, notifType string, path string) bool {
	var allowed bool = false
	for _, notif := range p.AllowedNotifs {
		if notif == notifType {
			allowed = true
		}
	}
	if allowed {
		message := pushover.NewMessageWithTitle(description, title)
		_, err := p.bot.SendMessage(message, p.recipient)
		if err != nil {
			p.Log.Error(err.Error())
			return false
		}
		return true
	}
	return false
}

func (p *Pushover) Connect() bool {
	p.bot = pushover.New(p.apiToken)
	p.recipient = pushover.NewRecipient(p.config.GetString("recipient"))
	if p.bot != nil {
		p.Log.Info("Connected to pushover")
		return true
	} else {
		p.Log.Warn("Failed to connect to pushover")
		return false
	}
}

func (p *Pushover) FromConfig(config viper.Viper) {
	p.config = config
	p.apiToken = config.GetString("apitoken")
	p.AllowedNotifs = config.GetStringSlice("notificationtypes")
}
