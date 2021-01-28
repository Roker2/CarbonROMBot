package commands

import (
	"carbonrombot/modules/carbonrom"
	"carbonrombot/modules/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	msgSupport = "CarbonROM supports %s."
	msgNotSupport = "CarbonROM doesn't support %s."
	getCarbonrom = "https://get.carbonrom.org/device-%s.html"
)

// alldevices
func AllDevices(ctx *ext.Context) error {
	devices, err := carbonrom.GetDevices()
	if err != nil {
		return err
	}
	msg := "Officially supported devices:"
	for _, device := range devices {
		msg += "\n• <code>" + device + "</code>"
	}
	_, err = ctx.EffectiveMessage.Reply(ctx.Bot, msg, &gotgbot.SendMessageOpts{ParseMode: "html"})
	return err
}

// device
func GetDevice(ctx *ext.Context) error {
	// If user didn't write device, bot should to say it
	// ctx.Args()[0] is a command
	if len(ctx.Args()) == 1 {
		_, err := ctx.EffectiveMessage.Reply(ctx.Bot, "You didn't write the device codename! Please write the device codename. Example:\n" +
			"<code>/device mido</code>", &gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	// If user wrote device, go to process it
	devices, err := carbonrom.GetDevices()
	if err != nil {
		return err
	}
	// ctx.Args()[0] is a command
	if utils.ContainsString(devices, ctx.Args()[1]) {
		roms, err := carbonrom.GetDeviceRoms(ctx.Args()[1])
		if err != nil {
			return err
		}
		//log.Print(files)
		_, err = ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf(msgSupport, ctx.Args()[1]), &gotgbot.SendMessageOpts{
			ParseMode: "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{Text: "Download the latest build", Url: roms[len(roms) - 1].RomUrl()},
					},
					{
						{Text: "All builds", Url: fmt.Sprintf(getCarbonrom, ctx.Args()[1])},
					},
				},
			},
		})
	} else {
		_, err = ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf(msgNotSupport, ctx.Args()[1]), nil)
	}
	return err
}