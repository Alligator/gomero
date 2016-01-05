package plugin

import (
	"log"
	"net/http"
	"net/url"
)

/*
API Description

POST /send - send a message
	key		= your api key
	message = the message
	target	= where the message should go (e.g. #channel or SomeUser)
	private = don't log the message (will only log that the api was called)

POST /register - register a command
	key		= your api key
	name	= the command name
	cb		= the callback url

	the callback url will get POSTed to with these params:
		source	= the source of the message (#channel or SomeUser)
		command = the command that was used
		message = the message
*/

type HttpServerPlugin struct {
	CallbackCommands map[string]string
}

func (serv *HttpServerPlugin) Init(bot Bot) {
	serv.CallbackCommands = make(map[string]string)

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		key := r.PostFormValue("key")

		if user, ok := bot.Config.ApiKeys[key]; ok {
			msg := r.PostFormValue("message")
			target := r.PostFormValue("target")
			private := r.PostFormValue("private")

			if private != "true" {
				log.Printf("API [%s]: %s %s", user, target, msg)
			}

			if len(target) > 0 {
				bot.Conn.SendRaw("PRIVMSG " + target + " " + msg)
			} else {
				bot.Conn.SendRaw("PRIVMSG #sa-minecraft " + msg)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid API key\n"))
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		key := r.PostFormValue("key")

		if user, ok := bot.Config.ApiKeys[key]; ok {
			name := r.PostFormValue("name")
			cb := r.PostFormValue("cb")

			if len(name) == 0 || len(cb) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid request"))
				return
			}

			log.Printf("API [%s]: registered command '%s' at callback '%s'", user, name, cb)

			serv.CallbackCommands[name] = cb

		} else {
			w.Write([]byte("invalid API key"))
		}
	})

	http.ListenAndServe("0.0.0.0:42069", nil)
	log.Printf("http server listening on 0.0.0.0:42069")
}

func (serv *HttpServerPlugin) Message(name string, ctx Context, bot Bot) {
	if cb, ok := serv.CallbackCommands[name]; ok {
		data := url.Values{}
		data.Add("source", ctx.Nick)
		data.Add("command", name)
		data.Add("message", ctx.Text)

		log.Printf("API POSTing %v to %s", data, cb)
		_, err := http.PostForm(cb, data)

		if err != nil {
			log.Printf("API POST error %s", err.Error())
		}
	}
}
