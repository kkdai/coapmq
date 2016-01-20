package coapmq

import (
	"log"
	"net"

	"github.com/dustin/go-coap"
)

type chanMapStringList map[*net.UDPAddr][]string
type stringMapChanList map[string][]*net.UDPAddr

type Broker struct {
	//The default maximal map size
	Capacity int

	msgIndex uint16 //for increase and sync message ID

	//map to store "chan -> Topic List" for find subscription
	clientMapTopics chanMapStringList
	//map to store "topic -> chan List" for publish
	topicMapClients stringMapChanList
	//Store all topic list and its latest value
	topicMapValue map[string]string
}

//Create a new pubsub server using CoAP protocol
//maxChannel: It is the subpub topic limitation size, suggest not lower than 1024 for basic usage
func NewBroker(maxCapacity int) *Broker {
	cSev := new(Broker)
	cSev.Capacity = maxCapacity
	cSev.clientMapTopics = make(map[*net.UDPAddr][]string, maxCapacity)
	cSev.topicMapClients = make(map[string][]*net.UDPAddr, maxCapacity)
	cSev.topicMapValue = make(map[string]string, maxCapacity)

	cSev.msgIndex = GetIPv4Int16() + GetLocalRandomInt()
	log.Println("Init msgID=", cSev.msgIndex)
	return cSev
}

func (c *Broker) getMsgID() uint16 {
	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *Broker) removeSubscription(topic string, client *net.UDPAddr) coap.COAPCode {
	res := coap.Deleted
	removeIndexT2C := -1
	if val, exist := c.topicMapClients[topic]; exist {
		for k, v := range val {
			if v == client {
				removeIndexT2C = k
			}
		}
		if removeIndexT2C != -1 {
			sliceClients := c.topicMapClients[topic]
			if len(sliceClients) > 1 {
				c.topicMapClients[topic] = append(sliceClients[:removeIndexT2C], sliceClients[removeIndexT2C+1:]...)
			} else {
				delete(c.topicMapClients, topic)
			}
		}
	}

	removeIndexC2T := -1
	if val, exist := c.clientMapTopics[client]; exist {
		for k, v := range val {
			if v == topic {
				removeIndexC2T = k
			}
		}
		if removeIndexC2T != -1 {
			sliceTopics := c.clientMapTopics[client]
			if len(sliceTopics) > 1 {
				c.clientMapTopics[client] = append(sliceTopics[:removeIndexC2T], sliceTopics[removeIndexC2T+1:]...)
			} else {
				delete(c.clientMapTopics, client)
			}
		}
	}
	return res
}

//Create new topic in coapmq broker
func (c *Broker) createTopic(topic string) coap.COAPCode {
	res := coap.Created
	if _, exist := c.topicMapValue[topic]; exist {
		res = coap.Forbidden
		log.Println("Create topic failed, topic exist.")
		return res
	}

	c.topicMapValue[topic] = "" //default value for creation
	return res
}

func (c *Broker) subscribeTopic(topic string, client *net.UDPAddr) coap.COAPCode {
	res := coap.Created

	return res
}

func (c *Broker) addSubscription(topic string, client *net.UDPAddr) coap.COAPCode {
	res := coap.Created

	topicFound := false
	if val, exist := c.topicMapClients[topic]; exist {
		for _, v := range val {
			if v == client {
				topicFound = true
			}
		}
	}
	if topicFound == false {
		res = coap.NotFound
		c.topicMapClients[topic] = append(c.topicMapClients[topic], client)
	}

	clientFound := false
	if val, exist := c.clientMapTopics[client]; exist {
		for _, v := range val {
			if v == topic {
				clientFound = true
			}
		}
	}

	if clientFound == false {
		c.clientMapTopics[client] = append(c.clientMapTopics[client], topic)
	}
	return res
}

func (c *Broker) readTopic(topic string) (string, coap.COAPCode) {
	res := coap.Content

	var retValue string
	if value, exist := c.topicMapValue[topic]; !exist {
		res = coap.NotFound
		retValue = ""
	} else {
		retValue = value
	}

	log.Println("read finished")
	return retValue, res
}

func (c *Broker) publish(l *net.UDPConn, topic string, value string) coap.COAPCode {
	res := coap.Changed
	if clients, exist := c.topicMapClients[topic]; !exist {
		return coap.NotFound
	} else { //topic exist, publish it
		for _, client := range clients {
			c.publishMsg(l, client, topic, value)
			log.Println("topic->", topic, " PUB to ", client, " msg=", value)
		}
	}

	c.topicMapValue[topic] = value
	log.Println("pub finished")
	return res
}

func (c *Broker) handleCoAPMessage(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	cmd, err := MessageDecode(m)
	if err != nil {
		log.Println("Message decode err:", err)
		m.Code = coap.BadRequest
		return m
	}

	log.Println("cmd=", cmd)

	res := coap.BadRequest
	retValue := ""
	reqCmd := ""

	switch cmd.Type {
	case CMD_SUBSCRIBE:
		res = c.addSubscription(cmd.Topic, a)
		reqCmd = "Subscription:"
	case CMD_UNSUBSCRIBE:
		res = c.removeSubscription(cmd.Topic, a)
		reqCmd = "Reqmove Sub:"
	case CMD_PUBLISH:
		res = c.publish(l, cmd.Topic, string(m.Payload))
		reqCmd = "Publish:"
	case CMD_HEARTBEAT:
		m.Code = coap.Content
		reqCmd = "Heart Beat:"
	case CMD_CREATE:
		res = c.createTopic(cmd.Topic)
		reqCmd = "Create topic:"
	case CMD_READ:
		retValue, res = c.readTopic(cmd.Topic)
		reqCmd = "Read topic:"
	default:
		reqCmd = "Invalid Command:"
	}

	log.Println("Got cmd=", reqCmd, " from:", a)

	for k, v := range c.topicMapClients {
		log.Println("Topic=", k, " sub by client=>", v)
	}

	//Prepare response message
	return c.response(res, retValue, m)
}

//Start to listen udp port and serve request, until faltal eror occur
func (c *Broker) ListenAndServe(udpPort string) {
	log.Fatal(coap.ListenAndServe("udp", udpPort,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			return c.handleCoAPMessage(l, a, m)
		})))
}

func (c *Broker) response(res coap.COAPCode, data string, m *coap.Message) *coap.Message {
	m.Type = coap.Acknowledgement
	m.Code = res
	if data != "" {
		m.Payload = []byte(data)
	}
	return m
}

func (c *Broker) publishMsg(l *net.UDPConn, a *net.UDPAddr, topic string, msg string) {
	m := EncodeMessage(c.getMsgID(), CMD_PUBLISH, msg, topic)
	log.Printf("Transmitting %v msg=%s", m, msg)
	err := coap.Transmit(l, a, *m)
	if err != nil {
		log.Printf("Error on transmitter, stopping: %v", err)
		return
	}
}
