package coapmq

import (
	"errors"
	"log"
	"time"

	"github.com/dustin/go-coap"
)

type subConnection struct {
	channel   chan string
	clientCon *coap.Conn
}

type Client struct {
	msgIndex uint16
	serAddr  string
	subList  map[string]subConnection
}

// Create a pubsub client for CoAP protocol
// It will connect to server and make sure it alive and start heart beat
// To keep udp port open, we will send heart beat event to server every minutes
func NewClient(servAddr string) *Client {
	c := new(Client)
	c.subList = make(map[string]subConnection, 0)
	c.serAddr = servAddr

	//Connection check if any error
	_, err := c.sendReq(CMD_HEARTBEAT, "", "")
	if err != nil {
		log.Println("Cannot connect to server")
		return nil
	}
	//Start heart beat
	c.msgIndex = GetIPv4Int16() + GetLocalRandomInt()
	log.Println("Init msgID=", c.msgIndex)
	go c.heartBeat()
	return c
}

func (c *Client) Publish(topic string, data string) error {
	_, err := c.sendReq(CMD_PUBLISH, topic, data)
	if err != nil {
		log.Println("pub error:", err)
	}

	return err
}

//Add Subscription on topic and return a channel for user to wait data
func (c *Client) Subscription(topic string) (chan string, error) {
	if val, exist := c.subList[topic]; exist {
		//if topic already exist in sub, return and not send to server
		return val.channel, nil
	}

	conn, err := c.sendWaitingReq(CMD_SUBSCRIBE, topic, "")
	if err != nil {
		return nil, err
	}

	subChan := make(chan string)
	go c.waitSubResponse(conn, subChan, topic)

	//Add client connection into member variable for heart beat
	clientConn := subConnection{channel: subChan, clientCon: conn}
	c.subList[topic] = clientConn
	return subChan, nil
}

//Create topic on server
func (c *Client) CreateTopic(topic string) error {
	ret, err := c.sendReq(CMD_CREATE, topic, "")
	log.Println("Result:", ret.Code)
	return ErrorWrapper(ret.Code, err)
}

//Remove topic on server
func (c *Client) RemoveTopic(topic string) error {
	ret, err := c.sendReq(CMD_REMOVE, topic, "")
	log.Println("Result:", ret.Code)
	return ErrorWrapper(ret.Code, err)
}

//Discovery and query with topic filter
func (c *Client) DiscoveryTopic(queryFilter string) (string, error) {
	//TODO. Need work detail on query filter on serer side
	return "", errors.New("NOT IMPLEMENT")
}

//Read topic most updated value from server, return error if topic not exist
func (c *Client) ReadTopic(topic string) (string, error) {
	ret, err := c.sendReq(CMD_READ, topic, "")
	if err != nil {
		return "", ErrorWrapper(ret.Code, err)
	}

	log.Println("Result:", ret.Code, " value=", string(ret.Payload))
	return string(ret.Payload), nil
}

//Remove Subscribetion on topic
func (c *Client) UnsubscribeTopic(topic string) error {
	if _, exist := c.subList[topic]; !exist {
		//if topic not in sub list, return and not send to server
		return errors.New("Not subscribe this topic before.")
	}

	ret, err := c.sendReq(CMD_UNSUBSCRIBE, topic, "")
	log.Println("Result:", ret.Code, " Err=", err, " detail err=", ErrorWrapper(ret.Code, err))
	return ErrorWrapper(ret.Code, err)
}

func (c *Client) sendWaitingReq(cmd CMD_TYPE, topic string, msg string) (*coap.Conn, error) {
	reqMsg := EncodeMessage(c.getMsgID(), cmd, msg, topic)
	log.Println("path=", reqMsg.Path())
	conn, err := coap.Dial("udp", c.serAddr)
	if err != nil {
		log.Printf(">>Error dialing: %v \n", err)
		return nil, errors.New("Dial failed")
	}
	conn.Send(*reqMsg)
	log.Println("msg->", *reqMsg)
	return conn, err
}

func (c *Client) sendReq(cmd CMD_TYPE, topic string, msg string) (*coap.Message, error) {
	reqMsg := EncodeMessage(c.getMsgID(), cmd, msg, topic)
	log.Println("path=", reqMsg.Path())
	conn, err := coap.Dial("udp", c.serAddr)
	if err != nil {
		log.Printf(">>Error dialing: %v \n", err)
		return nil, errors.New("Dial failed")
	}
	return conn.Send(*reqMsg)
}

func (c *Client) waitSubResponse(conn *coap.Conn, ch chan string, topic string) {
	log.Println("start to wait sub")
	var rv *coap.Message
	var err error
	var keepLoop bool
	keepLoop = true
	for keepLoop {
		if rv != nil {
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			log.Printf("Got %s", rv.Payload)
		}
		rv, err = conn.Receive()

		if err == nil {
			ch <- string(rv.Payload)
		}

		time.Sleep(time.Second)
		if _, exist := c.subList[topic]; !exist {
			//sub topic already remove, leave loop
			log.Println("Loop topic:", topic, " already remove leave loop")
			keepLoop = false
		}
	}
	log.Println("Leave wait sub")
}

func (c *Client) getMsgID() uint16 {
	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *Client) heartBeat() {
	log.Println("Starting heart beat loop call")

	for {
		_, err := c.sendReq(CMD_HEARTBEAT, "", "")
		if err != nil {
			log.Fatal("Server lost!")
			return
		}
		log.Println("Send the heart beat")
		time.Sleep(time.Minute)
	}
}
