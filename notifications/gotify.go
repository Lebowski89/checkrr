package notifications

import (
	"fmt"
	"github.com/aetaric/checkrr/logging"
	"net/http"
	"net/url"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/spf13/viper"
)

type GotifyNotifs struct {
	URL           string
	Client        *client.GotifyREST
	AuthToken     string
	Connected     bool
	AllowedNotifs []string
	Log           *logging.Log
}

func (d *GotifyNotifs) FromConfig(config viper.Viper) {
	d.URL = config.GetString("URL")
	d.AllowedNotifs = config.GetStringSlice("notificationtypes")
	d.AuthToken = config.GetString("authtoken")
}

func (d *GotifyNotifs) Connect() bool {
	myURL, _ := url.Parse(d.URL)
	client := gotify.NewClient(myURL, &http.Client{})
	versionResponse, err := client.Version.GetVersion(nil)

	if err != nil {
		d.Log.Warn("unable to connect to gotify")
		return false
	}
	version := versionResponse.Payload
	d.Client = client
	d.Connected = true
	d.Log.Info(fmt.Sprintf("Connected to Gotify, %s", version))
	return true
}

func (d GotifyNotifs) Notify(title string, description string, notifType string, path string) bool {
	if d.Connected {
		var allowed bool = false
		for _, notif := range d.AllowedNotifs {
			if notif == notifType {
				allowed = true
			}
		}
		if allowed {
			params := message.NewCreateMessageParams()
			params.Body = &models.MessageExternal{
				Title:    title,
				Message:  description,
				Priority: 5,
			}
			_, err := d.Client.Message.CreateMessage(params, auth.TokenAuth(d.AuthToken))

			return err == nil
		}
	}
	return false
}
