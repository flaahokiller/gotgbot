package main

import (
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"))
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(b, nil)
	dispatcher := updater.Dispatcher

	// Add echo handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("source", source))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("start_callback"), startCB))
	dispatcher.AddHandler(handlers.NewMessage(filters.All, echo))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{Clean: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func source(ctx *ext.Context) error {
	f, err := os.Open("samples/echoBot/main.go")
	if err != nil {
		fmt.Println("failed to open source: " + err.Error())
		return nil
	}

	_, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, f, &gotgbot.SendDocumentOpts{
		Caption:          "Here is my source code.",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	})
	if err != nil {
		fmt.Println("failed to send source: " + err.Error())
		return nil
	}

	// Alternative file sending solutions:

	// --- By file_id:
	//_, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, "file_id", &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	//})
	//if err != nil {
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	//}

	// --- By []byte:
	//bs, err := ioutil.ReadFile("samples/echoBot/main.go")
	//if err != nil {
	//	fmt.Println("failed to open source: " + err.Error())
	//	return nil
	//}
	//
	//_, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, bs, &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	//})
	//if err != nil {
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	//}

	// --- By custom name:
	//f2, err := os.Open("samples/echoBot/main.go")
	//if err != nil {
	//	fmt.Println("failed to open source: " + err.Error())
	//	return nil
	//}
	//
	//_, err = ctx.Bot.SendDocument(ctx.EffectiveChat.Id, gotgbot.NamedFile{
	//	File:     f2,
	//	FileName: "NewFileName",
	//}, &gotgbot.SendDocumentOpts{
	//	Caption:          "Here is my source code.",
	//	ReplyToMessageId: ctx.EffectiveMessage.MessageId,
	//})
	//if err != nil {
	//	fmt.Println("failed to send source: " + err.Error())
	//	return nil
	//}

	return nil
}

// start introduces the bot
func start(ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf("Hello, I'm @%s. I <b>repeat</b> all your messages.", ctx.Bot.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Press me", CallbackData: "start_callback"},
			}},
		},
	})
	if err != nil {
		fmt.Println("failed to send: " + err.Error())
	}
	return nil
}

// startCB edits the start message
func startCB(ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cb.Answer(ctx.Bot, nil)
	cb.Message.EditText(ctx.Bot, "You edited the start message.", nil)
	return nil
}

// echo replies to a messages with its own contents
func echo(ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(ctx.Bot, ctx.EffectiveMessage.Text, nil)
	return nil
}
