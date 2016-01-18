package coapmq

import (
	"strings"

	"github.com/dustin/go-coap"
)

type Cmd struct {
	Type   CMD_TYPE
	Topics []string
	Msg    string
}

func GetMsgCmdCode(cmd CMD_TYPE) coap.COAPCode {
	var code coap.COAPCode

	switch cmd {
	case CMD_DISCOVER:
		code = coap.GET
	case CMD_CREATE:
		code = coap.POST
	case CMD_PUBLISH:
		code = coap.PUT
	case CMD_SUBSCRIBE:
		code = coap.GET
	case CMD_UNSUBSCRIBE:
		code = coap.GET
	case CMD_READ:
		code = coap.GET
	case CMD_REMOVE:
		code = coap.DELETE
	}

	return code
}

//Refer to coapmq RFC:  https://datatracker.ietf.org/doc/draft-koster-core-coap-pubsub
//URI Template:  /{+ps/}{topic}{/topic*}
func EncodeCmdsToPath(cmd CMD_TYPE, topics []string) []string {
	var pathURI []string

	//add basic command here
	// cmd URI: {+ps}/{*topic}
	pathURI = append(pathURI, "ps")

	for _, v := range topics {
		pathURI = append(pathURI, v)
	}

	if cmd == CMD_DISCOVER {
		//TODO add query filter
	}
	return pathURI
}

func EncodeMessage(msgID uint16, cmd CMD_TYPE, msg string, topics ...string) *coap.Message {
	m := new(coap.Message)
	m.Type = coap.Confirmable
	m.Code = GetMsgCmdCode(cmd)
	m.MessageID = msgID

	m.Payload = []byte(msg)
	m.SetPath(EncodeCmdsToPath(cmd, topics))

	//m.SetOption(coap.ContentFormat, coap.TextPlain)
	m.SetOption(coap.ContentFormat, coap.AppLinkFormat)

	//specific handle for Observe (Refer RFC 7461)
	switch cmd {
	case CMD_SUBSCRIBE:
		m.SetOption(coap.Observe, 0)
	case CMD_UNSUBSCRIBE:
		m.SetOption(coap.Observe, 1)
	}
	return m
}

//Parse receive message to Coapmq.Cmd to get command and topic
func MessageDecode(m *coap.Message) *Cmd {
	path := m.Path()
	if len(path) == 0 || path[0] != "ps" {
		//cmd is not valid.
		return nil
	}

	var topic string
	if len(path) > 1 {
		topic = path[1]
	}

	c := new(Cmd)
	c.Type = CMD_INVALID
	c.Topics = append(c.Topics, topic)

	switch m.Code {
	case coap.GET:
		obseve, found := m.Option(coap.Observe).(uint8)
		if found {
			switch obseve {
			case 0:
				c.Type = CMD_SUBSCRIBE
			case 1:
				c.Type = CMD_UNSUBSCRIBE
			}
		} else {
			if strings.HasPrefix(topic, "?") {
				//it is discover
				c.Type = CMD_DISCOVER
			} else if topic != "" {
				c.Type = CMD_READ
			} else {
				//cmd not valid
			}
		}
	case coap.POST:
		c.Type = CMD_CREATE
	case coap.PUT:
		c.Type = CMD_PUBLISH
	case coap.DELETE:
		c.Type = CMD_REMOVE
	}

	if c.Type != CMD_INVALID {
		c.Msg = string(m.Payload)
	}
	return c
}
