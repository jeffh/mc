package mc

import (
	"ax"
	"io"
	"mc/protocol"
	//"unicode/utf16"
)

//////////////////////////////////////////////////////////
type Client struct {
	Connection    *protocol.Connection
	Outbox        chan interface{}
	Inbox         chan interface{}
	Logger        ax.Logger
	Exit          chan bool
	LogTraffic    bool
	AutoKeepAlive bool
}

func NewClient(stream io.ReadWriteCloser, msgBuffer int, l ax.Logger) *Client {
	reader := protocol.NewReader(stream, protocol.ClientPacketMapper, nil, l)
	writer := protocol.NewWriter(stream, protocol.ClientPacketMapper, nil, l)
	return &Client{
		Connection:    protocol.NewConnection(reader, writer),
		Outbox:        make(chan interface{}, msgBuffer),
		Inbox:         make(chan interface{}, msgBuffer),
		Logger:        ax.Wrap(ax.Use(l), ax.NewPrefixLogger("[client] ")),
		Exit:          make(chan bool),
		AutoKeepAlive: true,
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
			c.Logger.Printf("[mc/Client] Logging in with encryption %v", handshake)
		} else {
			c.Logger.Printf("[mc/Client] Logging in with NO encryption %v", handshake)
		}
	}

	if useEncryption {
		secret, err := protocol.GenerateSecretKey()
		if err != nil {
			return err
		}

		err = protocol.EstablishEncryptedConnection(c.Connection, handshake, secret)
		if err != nil {
			c.Logger.Printf("Failed to connect: %s", err)
			return err
		}
		protocol.EncryptConnection(c.Connection)
	} else {
		err = protocol.EstablishPlaintextConnection(c.Connection, handshake)
		if err != nil {
			c.Logger.Printf("Failed to connect: %s", err)
			return err
		}
	}

	// spawn!
	return err
}

func (c *Client) ConnectUnencrypted(hostname string, port int32, username string) error {
	return c.performConnect(hostname, port, username, false)
}

func (c *Client) ConnectEncrypted(hostname string, port int32, username string) error {
	return c.performConnect(hostname, port, username, true)
}

func (c *Client) ProcessInbox() {
	for {
		p, err := c.ReadPacket()
		if err == nil {
			c.Inbox <- p
		} else {
			c.Logger.Printf("Failed to read packet: %s", err)
			panic(err)
			if err == io.EOF {
				c.Exit <- true
				return
			}
		}
	}
}

func (c *Client) ProcessOutbox() {
	for {
		p, ok := <-c.Outbox
		if !ok {
			c.Logger.Printf("Outbox closed")
			c.Exit <- true
			return
		}
		err := c.WritePacket(p)
		if err != nil {
			c.Logger.Printf("Failed to write struct (%#v): %s", p, err)
			if err == io.EOF {
				c.Exit <- true
				return
			}
		}
	}
}

func (c *Client) WritePacket(v interface{}) error {
	err := c.Connection.WritePacket(v)
	return err
}

func (c *Client) ReadPacket() (interface{}, error) {
	p, err := c.Connection.ReadPacket()
	return p, err
}
