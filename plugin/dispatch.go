package plugin

import (
	"github.com/alligator/gomero/config"
	"github.com/alligator/gomero/db"
	"github.com/alligator/gomero/irc"
	ircLib "github.com/sorcix/irc"
	"strings"
)

type Dispatcher struct {
	PM  *PluginManager
	Bot Bot
}

type Bot struct {
	Conn   *irc.IrcConn
	Config config.Config
	Db     *db.Db
}

// uugh
func NewContext(msg ircLib.Message, cmd string, bot Bot) Context {
	ctx := Context{}
	ctx.Message = msg

	if cmd != "" {
		s := strings.SplitN(msg.Trailing, " ", 2)
		if len(s) > 1 {
			ctx.Text = s[1]
		}
	} else {
		ctx.Text = msg.Trailing
	}

	ctx.Nick = msg.Prefix.Name
	ctx.Host = msg.Prefix.Host
	ctx.Bot = bot

	if len(msg.Params) > 0 {
		if strings.HasPrefix(msg.Params[0], "#") {
			ctx.Channel = msg.Params[0]
		} else {
			ctx.Channel = ctx.Nick
		}
	}

	return ctx
}

func NewDispatcher(conn *irc.IrcConn, config config.Config, db *db.Db) *Dispatcher {
	d := new(Dispatcher)
	d.Bot = Bot{conn, config, db}
	d.PM = NewPluginManager("plugin/lua", d.Bot)
	go d.readLoop()
	return d
}

func (d *Dispatcher) readLoop() {
	for msg := range d.Bot.Conn.Out {
		go d.dispatch(msg)
	}
}

func (d *Dispatcher) dispatch(msg ircLib.Message) {
	if strings.HasPrefix(msg.Trailing, d.Bot.Config.Prefix) {
		// command
		raw := strings.Trim(strings.TrimPrefix(msg.Trailing, d.Bot.Config.Prefix), " ")
		sp := strings.SplitN(raw, " ", 2)
		cmd := sp[0]
		context := NewContext(msg, cmd, d.Bot)
		resp, _ := d.PM.CallCommand(cmd, context, d.Bot)
		if resp != "" {
			d.Bot.Conn.Inp <- ircLib.Message{
				Command:  "PRIVMSG",
				Params:   msg.Params,
				Trailing: resp,
			}
		}
	} else {
		// event
		context := NewContext(msg, "", d.Bot)
		respChan, err := d.PM.CallEvent(msg.Command, context, d.Bot)
		if err == nil {
			for resp := range respChan {
				d.Bot.Conn.SendRaw(resp)
			}
		}
	}
}
