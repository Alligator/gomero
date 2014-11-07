package plugin

import (
	"fmt"
	"github.com/alligator/gomero/config"
	"github.com/alligator/gomero/irc"
	ircLib "github.com/sorcix/irc"
	"strings"
)

type Dispatcher struct {
	Conn   *irc.IrcConn
	Config config.Config
	PM     *PluginManager
}

func NewDispatcher(conn *irc.IrcConn, config config.Config) *Dispatcher {
	d := new(Dispatcher)
	d.Conn = conn
	d.Config = config
	d.PM = NewPluginManager("plugin/lua")
	go d.readLoop()
	return d
}

func (d *Dispatcher) readLoop() {
	for msg := range d.Conn.Out {
		cmd := msg.Trailing
		if strings.HasPrefix(cmd, d.Config.Prefix) {
			go d.dispatch(msg)
		}
	}
}

func (d *Dispatcher) dispatch(msg ircLib.Message) {
	cmd := strings.TrimPrefix(msg.Trailing, d.Config.Prefix)
	resp, err := d.PM.Call(cmd)
	if err != nil {
		fmt.Print("DISPATCH")
		fmt.Println(err)
	} else {
		d.Conn.Inp <- ircLib.Message{
			Command:  "PRIVMSG",
			Params:   msg.Params,
			Trailing: resp,
		}
	}
}
