package plugin

import (
	ircLib "github.com/sorcix/irc"
)

type Context struct {
	Nick    string
	Host    string
	Channel string
	Text    string
	Message ircLib.Message
	Bot     Bot
}

func (ctx Context) Say(message string) {
	ctx.Bot.Conn.Inp <- ircLib.Message{
		Command:  "PRIVMSG",
		Params:   ctx.Message.Params,
		Trailing: message,
	}
}

func (ctx Context) Reply(message string) {
	message = ctx.Nick + ": " + message
	ctx.Bot.Conn.Inp <- ircLib.Message{
		Command:  "PRIVMSG",
		Params:   ctx.Message.Params,
		Trailing: message,
	}
}

func (ctx Context) Raw(message string) {
	ctx.Bot.Conn.SendRaw(message)
}
