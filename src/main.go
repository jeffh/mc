package main

import (
	"fmt"
	"mc"
	"mc/protocol"
	"net"
	"time"
)

func panicIfError(err error) {
	if err != nil {
		panic(fmt.Errorf("[Client] Error: %s\n", err))
	}
}

func main() {
	host := "localhost"
	port := int32(25565)
	hostport := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		panic(err)
	}
	c := mc.NewClient(conn, 100, &mc.StdoutLogger{})
	c.LogTraffic = true
	err = c.ConnectUnencrypted(host, port, "MCBot")
	//err = c.Connect("localhost", 1337, "MCBot")
	panicIfError(err)

	go c.ProcessInbox()
	go c.ProcessOutbox()

	go func() {
		var position *protocol.PlayerPositionLookForClient
		c.Outbox <- &protocol.ClientStatus{}

		for {
			timeout := time.After(50 * time.Millisecond)
			select {
			case pck := <-c.Inbox:
				switch t := pck.(type) {
				case *protocol.KeepAlive:
					c.Outbox <- t
				case *protocol.PlayerPositionLookForClient:
					position = t
					c.Outbox <- t.PacketForServer()
				case *protocol.Disconnect:
					c.Exit <- true
					return
				}
			case <-timeout:
				if position != nil {
					c.Outbox <- position.PacketForServer()
				}
			}
		}
	}()

	<-c.Exit
}
