package main

import (
	"carbonrombot/modules/commands"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"), &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
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
	dispatcher.AddHandler(handlers.NewCommand("device", commands.GetDevice))
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
		webhook := ext.WebhookOpts{
			Listen: "0.0.0.0",
			Port: port,
			URLPath: b.Token,
		}
		err = updater.StartWebhook(b, webhook)
		if err != nil {
			panic("failed to start webhook: " + err.Error())
		}
		ok, err := b.SetWebhook(herokuUrl + b.Token, &gotgbot.SetWebhookOpts{MaxConnections: 40})
		if err != nil {
			panic("failed to start webhook: " + err.Error())
		}
		if !ok {
			panic("failed to set webhook, ok is false")
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