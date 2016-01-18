package coapmq

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/dustin/go-coap"
)

type CMD_TYPE int

//Refer to CoAP pub/sub Function Set (chapter 4)
const (
	CMD_DISCOVER    CMD_TYPE = iota
	CMD_CREATE      CMD_TYPE = iota
	CMD_PUBLISH     CMD_TYPE = iota
	CMD_SUBSCRIBE   CMD_TYPE = iota
	CMD_UNSUBSCRIBE CMD_TYPE = iota
	CMD_READ        CMD_TYPE = iota
	CMD_REMOVE      CMD_TYPE = iota
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Parse interface (which is []uint8) to string
func ParseUint8ToString(in interface{}) string {
	val, ok := in.([]uint8)
	if ok {
		return string(val)
	} else {
		return ""
	}
}

//According to RFC 7252, message ID need IP with a random number
//Get random number locally
func GetLocalRandomInt() uint16 {
	rand.Seed(time.Now().UnixNano())
	return uint16(rand.Intn(1000))
}

//According to RFC 7252, we need indicate message ID with sender IP or address + random number
//Get local IPV4 address to uint16 by <<8
func GetIPv4Int16() uint16 {

	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		log.Println("No network:", err)
		return 0
	}

	for _, i := range ifaces {
		if strings.Contains(i.Name, "en0") {
			addrs, err := i.Addrs()
			// handle err
			if err != nil {
				log.Println("No IP:", err)
				return 0
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// process IP address
				if ip[0] == 0 {
					//target XX.XX.XX.XX ipv4
					var myIP uint16
					myIP = uint16(ip[12])<<8 + uint16(ip[13])<<7 + uint16(ip[14])<<6 + uint16(ip[13])<<6
					return myIP
				}
			}
		}
	}

	return 0
}

//Refer to coapmq RFC:  https://datatracker.ietf.org/doc/draft-koster-core-coap-pubsub
//URI Template:  /{+ps/}{topic}{/topic*}
func EncodeCmdsToPath(cmd string, topics ...string) string {
	return ""
}

func GetCmdMessage(cmd CMD_TYPE, msg string, topics ...string) *coap.Message {
}
