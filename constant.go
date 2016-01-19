package coapmq

import "github.com/dustin/go-coap"

type CMD_TYPE int

//Refer to CoAP pub/sub Function Set (chapter 4)
const (
	CMD_INVALID     CMD_TYPE = iota
	CMD_DISCOVER    CMD_TYPE = iota
	CMD_CREATE      CMD_TYPE = iota
	CMD_PUBLISH     CMD_TYPE = iota
	CMD_SUBSCRIBE   CMD_TYPE = iota
	CMD_UNSUBSCRIBE CMD_TYPE = iota
	CMD_READ        CMD_TYPE = iota
	CMD_REMOVE      CMD_TYPE = iota
	//Propietary command to keep UDP connection alive
	CMD_HEARTBEAT CMD_TYPE = iota
)

var ErrorCodeMappingTable map[coap.COAPCode]string = map[coap.COAPCode]string{
	coap.Created:       "Created",
	coap.Deleted:       "Deleted",
	coap.Changed:       "Changed",
	coap.Content:       "Content",
	coap.BadRequest:    "Bad Request",
	coap.Unauthorized:  "Unauthorized",
	coap.NotFound:      "Not Found",
	coap.Forbidden:     "Forbidden",
	coap.NotAcceptable: "Not Acceptable",
}
