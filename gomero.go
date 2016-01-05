package main

import (
	"github.com/alligator/gomero/config"
	"github.com/alligator/gomero/db"
	"github.com/alligator/gomero/irc"
	"github.com/alligator/gomero/plugin"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)

	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Printf("config parse error: %s\n", err.Error())
		return
	}

	ircConn := irc.NewIrcConn(config.Host)
	db := db.NewDb()

	ircConn.Dial(config.Name)
	_ = plugin.NewDispatcher(ircConn, config, db)

	for {
		time.Sleep(1 * time.Second)
	}
}
