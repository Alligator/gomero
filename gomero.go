package main

import (
	"github.com/alligator/gomero/config"
	"github.com/alligator/gomero/db"
	"github.com/alligator/gomero/irc"
	"github.com/alligator/gomero/plugin"
	ircLib "github.com/sorcix/irc"
	"time"
)

func main() {
	config := config.ReadConfig("config.json")

	db := db.NewDb()
	db = db
	ircConn := irc.NewIrcConn(config.Host)

	ircConn.Dial()
	_ = plugin.NewDispatcher(ircConn, config)

	time.Sleep(8 * time.Second)
	ircConn.Inp <- ircLib.Message{
		Command:  "JOIN",
		Trailing: "#sa-minecraft",
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
