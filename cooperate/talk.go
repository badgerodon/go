package main

import (
	"encoding/gob"
	"log"
	"net"
	"sync"
)

type (
	Talk struct {
		send, recv chan Message
		lastid uint
		lock sync.Mutex
		calls map[uint]chan Reply
		streams map[uint]chan []byte
	}
	Message struct {
		Id uint
		Payload interface{}
	}
	Request struct {
		Name string
		Args []interface{}
	}
	Reply struct {
		Result interface{}
		Error error
	}
	StreamStart struct {}
	Stream []byte
	StreamStop struct {}
)

func init() {
	gob.Register(Request{})
	gob.Register(Reply{})
	gob.Register(Message{})
	gob.Register(Stream{})
	gob.Register(StreamStart{})
	gob.Register(StreamStop{})
}

func NewTalk(conn net.Conn) *Talk {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	send := make(chan Message, 1)
	recv := make(chan Message, 1)

	go func() {
		for {
			msg, ok := <- send
			if !ok {
				break
			}
			err := enc.Encode(msg)
			if err != nil {
				recv <- Message{0, err}
				break
			}
		}
	}()
	go func() {
		for {
			var msg Message
			err := dec.Decode(&msg)
			if err != nil {
				recv <- Message{0, err}
				break
			}

			recv <- msg
		}
	}()

	return &Talk{
		send: send,
		recv: recv,
		calls: make(map[uint]chan Reply),
		streams: make(map[uint]chan []byte),
	}
}

func (this *Talk) Send(req Request) Reply {
	this.lock.Lock()
	id := this.lastid
	this.lastid++
	reply := make(chan Reply, 1)
	this.calls[id] = reply
	this.lock.Unlock()

	this.send <- Message{id, req}

	return <- reply
}
func (this *Talk) Recv(handler func(Request)Reply) {
	for {
		msg := <- this.recv
		log.Println("MSG", msg)
		switch obj := msg.Payload.(type) {
		case error:
			return
		case StreamStart:
			this.lock.Lock()
			stream := make(chan []byte)
			this.streams[msg.Id] = stream
			c, ok := this.calls[msg.Id]
			if ok {
				delete(this.calls, msg.Id)
			}
			this.lock.Unlock()

			if ok {
				c <- Reply{stream, nil}
			}
		case Stream:
			this.lock.Lock()
			stream, ok := this.streams[msg.Id]
			this.lock.Unlock()
			if ok {
				stream <- []byte(obj)
			}
		case StreamStop:
			this.lock.Lock()
			stream, ok := this.streams[msg.Id]
			if ok {
				delete(this.streams, msg.Id)
				close(stream)
			}
			this.lock.Unlock()
		case Request:
			id := msg.Id
			go func() {
				reply := handler(obj)
				c, ok := reply.Result.(chan []byte)
				if ok {
					this.send <- Message{id, StreamStart{}}
					for {
						bs, ok := <- c
						if !ok {
							this.send <- Message{id, StreamStop{}}
							break
						} else {
							this.send <- Message{id, Stream(bs)}
						}
					}
				} else {
					this.send <- Message{id, reply}
				}
			}()
		case Reply:
			this.lock.Lock()
			c, ok := this.calls[msg.Id]
			if ok {
				delete(this.calls, msg.Id)
			}
			this.lock.Unlock()

			if ok {
				c <- obj
			}
		}
	}
}