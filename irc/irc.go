package irc

import (
	"github.com/sorcix/irc"
	"log"
	"time"
)

type IrcConn struct {
	Host string
	Conn *irc.Conn
	Inp  chan irc.Message
	Out  chan irc.Message
}

func NewIrcConn(host string) *IrcConn {
	ic := new(IrcConn)
	ic.Host = host
	ic.Inp = make(chan irc.Message, 100)
	ic.Out = make(chan irc.Message, 100)
	return ic
}

func (ic *IrcConn) Dial(name string) {
	conn, err := irc.Dial(ic.Host)
	if err != nil {
		return
	}
	ic.Conn = conn

	go ic.srvRecv()
	go ic.srvSend()

	ic.Inp <- irc.Message{
		Command:  "NICK",
		Trailing: name,
	}
	ic.Inp <- irc.Message{
		Command:  "USER",
		Params:   []string{name, "3", "*"},
		Trailing: name,
	}
}

/*** private methods ***/
func (ic *IrcConn) srvRecv() {
	for {
		msg, err := ic.Conn.Decode()
		if err != nil {
			panic(err)
		}

		if msg.Command == "PING" {
			pongMsg := irc.Message{
				Command:  "PONG",
				Trailing: msg.Trailing,
			}
			ic.Inp <- pongMsg
		} else {
			log.Printf("%s\n", msg.String())
			ic.Out <- *msg
		}
	}
}

func (ic *IrcConn) srvSend() {
	burstLimiter := make(chan time.Time, 5)
	go func() {
		for t := range time.Tick(time.Second * 2) {
			burstLimiter <- t
		}
	}()

	for msg := range ic.Inp {
		<-burstLimiter
		ic.Conn.Encode(&msg)
		if msg.Command != "PONG" {
			log.Printf(">>> %s\n", msg.String())
		}
	}
}

func (ic *IrcConn) SendRaw(raw string) {
	ic.Inp <- *irc.ParseMessage(raw)
}
