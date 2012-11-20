package mc

import (
	"fmt"
	"io"
	"log"
	"mc/protocol"
	"os"
	//"unicode/utf16"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

//////////////////////////////////////////////////////////
type Client struct {
	Connection    *protocol.Connection
	Outbox        chan interface{}
	Inbox         chan interface{}
	Logger        Logger
	Exited        chan bool
	LogTraffic    bool
	AutoKeepAlive bool
}

func NewClient(stream io.ReadWriteCloser) *Client {
	reader := protocol.NewReader(stream, nil, protocol.ClientPacketMapper)
	writer := protocol.NewWriter(stream, nil, protocol.ClientPacketMapper)
	return &Client{
		Connection:    protocol.NewConnection(reader, writer),
		Outbox:        make(chan interface{}, 50),
		Inbox:         make(chan interface{}, 50),
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
		Exited:        make(chan bool),
		AutoKeepAlive: true,
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

//func (c *Client) readDump() {
//	for {
//		var b byte
//		err := c.Connection.ReadValue(&b)
//		fmt.Printf("-> 0x%x\n", b)
//		if err != nil {
//			panic(fmt.Errorf("Error: %s\n", err))
//		}
//	}
//}

func (c *Client) performConnect(hostname string, port int32, username string, useEncryption bool) (err error) {
	handshake := &protocol.Handshake{
		Version:  protocol.Version,
		Username: username,
		Hostname: hostname,
		Port:     port,
	}
	if c.LogTraffic {
		if useEncryption {
			c.log("[Logging in] (With Encryption) %v", handshake)
		} else {
			c.log("[Logging in] (No Encryption) %v", handshake)
		}
	}

	if useEncryption {
		secret, err := protocol.GenerateSecretKey()
		if err != nil {
			return err
		}

		err = protocol.EstablishEncryptedConnection(c.Connection, handshake, secret)
		if err != nil {
			c.log("Failed to connect: %s", err)
			return err
		}
		protocol.EncryptConnection(c.Connection)
	} else {
		err = protocol.EstablishPlaintextConnection(c.Connection, handshake)
		if err != nil {
			c.log("Failed to connect: %s", err)
			return err
		}
	}

	c.Outbox <- &protocol.ClientStatus{}

	// spawn!
	return err
}

func (c *Client) ConnectUnencrypted(hostname string, port int32, username string) error {
	return c.performConnect(hostname, port, username, false)
}

func (c *Client) Connect(hostname string, port int32, username string) error {
	return c.performConnect(hostname, port, username, true)
}

func (c *Client) ProcessInbox() {
	for {
		p, err := c.ReadPacket()
		if err == nil {
			c.Inbox <- p
		} else {
			c.log("Failed to read packet: %s", err)
			if err == io.EOF {
				c.Exited <- true
				return
			}
		}
	}
}

func (c *Client) ProcessOutbox() {
	for {
		p, ok := <-c.Outbox
		if !ok {
			c.log("Outbox closed")
			c.Exited <- true
			return
		}
		err := c.WritePacket(p)
		if err != nil {
			c.log("Failed to write struct (%#v): %s", p, err)
			if err == io.EOF {
				c.Exited <- true
				return
			}
		}
	}
}

func (c *Client) log(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	c.Logger.Printf("[Client] %s", s)
}

func (c *Client) WritePacket(v interface{}) error {
	if c.LogTraffic {
		c.log("<- %#v", v)
	}
	err := c.Connection.WritePacket(v)
	return err
}

func (c *Client) ReadPacket() (interface{}, error) {
	p, err := c.Connection.ReadPacket()
	if c.LogTraffic {
		c.log("-> %#v", p)
	}
	return p, err
}
