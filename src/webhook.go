package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Webhook struct {
	Name      string      `json:"name"`
	Type      int         `json:"type"`
	ChannelID string      `json:"channel_id"`
	Token     string      `json:"token"`
	Avatar    interface{} `json:"avatar"`
	GuildID   string      `json:"guild_id"`
	ID        string      `json:"id"`
	User      struct {
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		ID            string `json:"id"`
		Avatar        string `json:"avatar"`
	} `json:"user"`
}

// getWebhookMap provides a map of channelID string to Webhooks
func (conf Config) getWebhookMap(guilds []Guild) map[string]Webhook {
	var knownHooks = make(map[string]Webhook)
	for _, guild := range guilds {
		hooks := conf.getWebhooksForGuild(guild)
		for _, hook := range hooks {
			knownHooks[string(hook.ChannelID)] = hook
		}
	}
	return knownHooks
}

// getWebhooksForGuild returns all Webhooks for a Guild
func (conf Config) getWebhooksForGuild(guild Guild) []Webhook {
	// Build the appropriate request
	fullURL := conf.baseURL + "/guilds/" + string(guild.Id) + "/webhooks"
	request, _ := http.NewRequest("GET", fullURL, nil)
	request.Header.Add("Authorization", conf.token)

	// Make the request
	resp, err := conf.client.Do(request)
	if err != nil {
		log.Printf("Error making request %+v", err)
	}
	defer resp.Body.Close()

	// Parse the response into the data we return
	var webhooks []Webhook
	bodyJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("%+v", err)
	}
	err = json.Unmarshal([]byte(bodyJson), &webhooks)
	if err != nil {
		fmt.Errorf("%+v", err)
	}
	return webhooks

}
