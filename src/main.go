package main

import (
    "fmt"
    "net"
    "mc"
)
func main() {
    host := "localhost"
    port := int32(25565)
    hostport := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", hostport)
    if err != nil { panic(err) }
    c := mc.NewClient(conn)
    c.LogTraffic = true
    err = c.ConnectUnencrypted(host, port, "MCBot")
    //err = c.Connect("localhost", 1337, "MCBot")
    if err != nil {
        panic(fmt.Errorf("[Client] Error: %s\n", err))
    }
    go c.ProcessInbox()
    go c.ProcessOutbox()
    <-c.Exited
}
