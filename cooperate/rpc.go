package main

import (
	"encoding/gob"
	"errors"
	"log"
	"net"
	"sync"
)

type (
	RPC struct {
		c net.Conn
		enc *gob.Encoder
		dec *gob.Decoder
		lock sync.Mutex
		lastid uint
		calls map[uint]chan Reply
		handlers map[string]func(Request)Reply
	}

	Message struct {
		Type byte
		Id uint
		Payload interface{}
	}
	Request struct {
		Name string
		Args []interface{}
	}
	Reply struct {
		Data interface{}
		Error error
	}
)

const (
	TypeRequest byte = iota
	TypeReply
)

func init() {
	gob.Register(Request{})
	gob.Register(Reply{})
	gob.Register(Message{})
}

func NewRPC(conn net.Conn) *RPC {
	var m sync.Mutex
	return &RPC{
		conn,
		gob.NewEncoder(conn),
		gob.NewDecoder(conn),
		m,
		0,
		make(map[uint]chan Reply),
		make(map[string]func(Request)Reply),
	}
}
func (this *RPC) Call(name string, args ... interface{}) (interface{}, error) {
	output := make(chan Reply)

	this.lock.Lock()
	id := this.lastid
	this.lastid++
	this.calls[id] = output
	this.lock.Unlock()

	err := this.request(id, Request{name, args})
	if err != nil {
		this.lock.Lock()
		delete(this.calls, id)
		this.lock.Unlock()

		return nil, err
	}

	reply := <- output

	this.lock.Lock()
	delete(this.calls, id)
	this.lock.Unlock()

	return reply.Data, reply.Error
}
func (this *RPC) On(name string, handler func(req Request) Reply) {
	this.lock.Lock()
	this.handlers[name] = handler
	this.lock.Unlock()
}
func (this *RPC) Listen() error {
	for {
		var msg Message
		err := this.dec.Decode(&msg)
		log.Println("MESSAGE", msg)
		if err != nil {
			return err
		}
		switch msg.Type {
		case TypeRequest:
			req, ok := msg.Payload.(Request)
			if !ok {
				err = this.reply(msg.Id, Reply{nil,errors.New("Error Decoding")})
				if err != nil {
					return err
				}
				continue
			}

			this.lock.Lock()
			handler, ok := this.handlers[req.Name]
			this.lock.Unlock()

			if !ok {
				err = this.reply(msg.Id, Reply{nil,errors.New("Unknown Method")})
				if err != nil {
					return err
				}
				continue
			}

			err = this.reply(msg.Id, handler(req))
			if err != nil {
				return err
			}
		case TypeReply:
			reply, ok := msg.Payload.(Reply)
			if !ok {
				this.handleReply(msg.Id, Reply{nil,errors.New("Error Decoding")})
				continue
			}

			this.handleReply(msg.Id, reply)
		}
	}
	return nil
}
func (this *RPC) handleReply(id uint, reply Reply) {
	this.lock.Lock()
	c, ok := this.calls[id]
	this.lock.Unlock()
	if ok {
		c <- reply
	}
}
func (this *RPC) request(id uint, req Request) error {
	return this.enc.Encode(Message{
		Type: TypeRequest,
		Id: id,
		Payload: req,
	})
}
func (this *RPC) reply(id uint, reply Reply) error {
	return this.enc.Encode(Message{
		Type: TypeReply,
		Id: id,
		Payload: reply,
	})
}