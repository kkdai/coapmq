package coapmq

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
