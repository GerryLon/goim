package main

import (
	"goim/libs/bufio"
	"goim/libs/proto"
)

// Channel used by message pusher send msg to write goroutine.
// Channel维护一个长连接用户，只能对应一个 Room. 推送的消息可以在 Room 内广播，也可以推送到指定的 Channel.
type Channel struct {
	Room     *Room
	CliProto Ring
	signal   chan *proto.Proto
	Writer   bufio.Writer
	Reader   bufio.Reader
	Next     *Channel
	Prev     *Channel
}

func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *proto.Proto, svr)
	return c
}

// Push server push message.
func (c *Channel) Push(p *proto.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *proto.Proto {
	return <-c.signal
}

// Signal send signal to the channel, protocol ready.
func (c *Channel) Signal() {
	c.signal <- proto.ProtoReady
}

// Close close the channel.
func (c *Channel) Close() {
	c.signal <- proto.ProtoFinish
}
