package commands

import (
	"carbonrombot/modules/carbonrom"
	"carbonrombot/modules/utils"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	msgSupport    = "CarbonROM supports %s.\nThe latest ROM info:\n<b>ROM name:</b> %s\n<b>ROM version:</b> %s\n<b>MD5:</b> <code>%s</code>\n<b>Build date:</b> %s"
	msgNotSupport = "CarbonROM doesn't support %s."
	getCarbonrom  = "https://get.carbonrom.org/device-%s.html"
)

// start and help
func Help(b *gotgbot.Bot, ctx *ext.Context) error {
	opts := gotgbot.GetMyCommandsOpts{Scope: gotgbot.BotCommandScopeDefault{}}
	commands, err := b.GetMyCommands(&opts)
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		return fmt.Errorf("Commands list is empty")
	}
	msgText := "It's bot for getting info about supporting your device and getting the latest update.\n" +
		"Available commands:\n"
	for _, command := range commands {
		msgText += fmt.Sprintf("/%s - %s\n", command.Command, command.Description)
	}
	msgText += "Enjoy!"
	_, err = ctx.EffectiveMessage.Reply(b, msgText, nil)
	return err
}

// alldevices
func AllDevices(b *gotgbot.Bot, ctx *ext.Context) error {
	devices, err := carbonrom.GetDevices()
	if err != nil {
		return err
	}
	msg := "Officially supported devices:"
	for _, device := range devices {
		msg += "\n• <code>" + device + "</code>"
	}
	_, err = ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{ParseMode: "html"})
	return err
}

// device
func GetDevice(b *gotgbot.Bot, ctx *ext.Context) error {
	// If user didn't write device, bot should to say it
	// ctx.Args()[0] is a command
	if len(ctx.Args()) == 1 {
		_, err := ctx.EffectiveMessage.Reply(b, "You didn't write the device codename! Please write the device codename. Example:\n"+
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
		latestRom := roms[len(roms)-1]
		md5, err := latestRom.Md5()
		if err != nil {
			return err
		}
		androidVersion, err := latestRom.RomVersion()
		// I think showing the error in the logs it's okay
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf(msgSupport, ctx.Args()[1], latestRom.RomName(), androidVersion, md5, latestRom.GetDateAsString()), &gotgbot.SendMessageOpts{
			ParseMode: "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{Text: "Download the latest build", Url: latestRom.RomUrl()},
					},
					{
						{Text: "All builds", Url: fmt.Sprintf(getCarbonrom, ctx.Args()[1])},
					},
				},
			},
		})
	} else {
		_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf(msgNotSupport, ctx.Args()[1]), nil)
	}
	return err
}

func VersionsList(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := "ROM versions:"
	for romVersion, androidVersion := range carbonrom.Versions {
		msg += "\n• " + romVersion + " - " + androidVersion
	}
	_, err := ctx.EffectiveMessage.Reply(b, msg, nil)
	return err
}
