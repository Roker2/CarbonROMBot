package main

import (
	"carbonrombot/modules/commands"
	"fmt"
	"os"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	botToken := os.Getenv("TOKEN")
	if botToken == "" {
		panic("TOKEN environment variable is empty")
	}
	// Create bot from environment value.
	b, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Add handlers
	dispatcher.AddHandler(handlers.NewCommand("start", commands.Help))
	dispatcher.AddHandler(handlers.NewCommand("help", commands.Help))
	dispatcher.AddHandler(handlers.NewCommand("alldevices", commands.AllDevices))
	dispatcher.AddHandler(handlers.NewCommand("devices", commands.AllDevices))
	dispatcher.AddHandler(handlers.NewCommand("device", commands.GetDevice))
	dispatcher.AddHandler(handlers.NewCommand("romversions", commands.VersionsList))
	dispatcher.Error = func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
		fmt.Println(err.Error())
		return ext.DispatcherActionNoop
	}

	// Start receiving updates.
	if os.Getenv("USE_WEBHOOKS") == "yes" {
		fmt.Println("Starting webhook")
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic("failed to get port: " + err.Error())
		}
		herokuUrl := os.Getenv("HEROKU_URL")
		if herokuUrl == "" {
			panic("HEROKU_URL environment variable is empty")
		}
		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		if webhookSecret == "" {
			panic("WEBHOOK_SECRET environment variable is empty")
		}
		webhookOpts := ext.WebhookOpts{
			ListenAddr:  fmt.Sprintf("0.0.0.0:%d", port),
			SecretToken: webhookSecret,
		}
		err = updater.StartWebhook(b, botToken, webhookOpts)
		if err != nil {
			panic("failed to start webhook: " + err.Error())
		}
		err = updater.SetAllBotWebhooks(herokuUrl, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: true,
			SecretToken:        webhookOpts.SecretToken,
		})
		if err != nil {
			panic("failed to set webhook: " + err.Error())
		}
	} else {
		err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
		if err != nil {
			panic("failed to start polling: " + err.Error())
		}
	}

	fmt.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
