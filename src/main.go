package main

import (
	"ax"
	"fmt"
	"mc"
	"mc/protocol"
	"mc/simulator"
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
	var logger ax.Logger
	//logger = &ax.NullLogger{}
	logger = &ax.StdoutLogger{}
	logger = ax.Wrap(logger, ax.NewTimestampLogger(), ax.NewLockedLogger())
	c := mc.NewClient(conn, 20, logger)
	c.LogTraffic = true
	err = c.ConnectUnencrypted(host, port, "MCBot")
	//err = c.Connect("localhost", 1337, "MCBot")
	panicIfError(err)

	go c.ProcessInbox()
	go c.ProcessOutbox()

	go func() {
		sim := simulator.NewSimulator(logger)
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
				sim.ProcessMessage(pck)
			case <-timeout:
				if position != nil {
					c.Outbox <- position.PacketForServer()
				}
			}
		}
	}()

	<-c.Exit
}
