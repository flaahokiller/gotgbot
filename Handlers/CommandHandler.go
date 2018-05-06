package Handlers

import (
	"strings"
	"gotgbot/Ext"
	"gotgbot"
)

type Command struct {
	command  string
	response func(b Ext.Bot, u gotgbot.Update)
}

func NewCommand(command string, response func(b Ext.Bot, u gotgbot.Update)) Command {
	h := Command{}
	h.command = command
	h.response = response
	return h
}

func (h Command) HandleUpdate(update gotgbot.Update, d gotgbot.Dispatcher) {
	go h.response(d.Bot, update)
}

func (h Command) CheckUpdate(update gotgbot.Update) bool {
	return update.Message != nil && update.Message.Text != "" &&
		strings.Split(strings.Fields(update.Message.Text)[0], "@")[0] == "/" + h.command
}