package irc

import (
	"github.com/sorcix/irc"
	"log"
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

func (ic *IrcConn) Dial() {
	conn, err := irc.Dial(ic.Host)
	if err != nil {
		return
	}
	ic.Conn = conn

	go ic.srvRecv()
	go ic.srvSend()

	ic.Inp <- irc.Message{
		Command:  "NICK",
		Trailing: "gomero",
	}
	ic.Inp <- irc.Message{
		Command:  "USER",
		Params:   []string{"gomero", "3", "*"},
		Trailing: "gomero",
	}
}

/*** private methods ***/
func (ic *IrcConn) srvRecv() {
	for {
		msg, err := ic.Conn.Decode()
		if err != nil {
			panic(err)
		}

		log.Println(msg.String())
		if msg.Command == "PING" {
			pongMsg := irc.Message{
				Command:  "PONG",
				Trailing: msg.Trailing,
			}
			ic.Inp <- pongMsg
		} else {
			ic.Out <- *msg
		}
	}
}

func (ic *IrcConn) srvSend() {
	for msg := range ic.Inp {
		ic.Conn.Encode(&msg)
		log.Printf(">>> %s\n", msg.String())
	}
}
