package plugin

import (
	"log"
	"net/http"
)

type HttpServerPlugin struct{}

func (serv *HttpServerPlugin) init(bot Bot) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.PostFormValue("key")

		if user, ok := bot.Config.ApiKeys[key]; ok {
			msg := r.PostFormValue("message")
			log.Printf("API [%s]: %s", user, msg)
			bot.Conn.SendRaw("PRIVMSG #sa-minecraft " + msg)
		}
	})

	http.ListenAndServe("0.0.0.0:42069", nil)
	log.Printf("http server listening on 0.0.0.0:42069")
}
