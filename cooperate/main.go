package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var (
	syncer *Syncer = NewSyncer()
)

func mkRequest(name string, args ... interface{}) Request {
	return Request{name, args}
}

func dosync(talk *Talk, folder string) error {
	defer talk.Send(Request{"close", nil})

	local, err := syncer.GetDigest(folder)
	if err != nil {
		return err
	}

	reply := talk.Send(mkRequest("get-digest", folder))
	if reply.Error != nil {
		return reply.Error
	}
	remote := reply.Result.(*Digest)

	toget := []string{}
	//todel := []string{}

	i := 0
	j := 0
	for {
		if !(i < len(local.Entries) || j < len(remote.Entries)) {
			break
		}
		// Not in local, get file
		if i >= len(local.Entries) {
			toget = append(toget, remote.Entries[j].Path)
			j++
			continue
		}
		// Not in remote
		if j >= len(remote.Entries) {
			i++
			continue
		}
		ei := local.Entries[i]
		ej := remote.Entries[j]
		// File in remote, not in local
		if ei.Path > ej.Path {
			toget = append(toget, remote.Entries[j].Path)
			j++
			continue
		}
		// File in local, not in remote
		if ei.Path < ej.Path {
			i++
			continue
		}

		// File in both
		// Different?
		if ei.Hash != ej.Hash {
			if ei.Time < ej.Time {
				toget = append(toget, remote.Entries[j].Path)
			}
		}
		i++
		j++
	}

	for _, p := range toget {
		reply = talk.Send(mkRequest("get-file", folder, p))
		if reply.Error != nil {
			return reply.Error
		}
		stream := reply.Result.(chan []byte)
		err = syncer.WriteFile(folder, p, stream)
		if err != nil {
			return err
		}
	}

	return nil
}

func listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

func recv(conn net.Conn, talk *Talk, waiter *Waiter) {
	talk.Recv(func(req Request) Reply {
		log.Println("RECV", req)
		switch req.Name {
		case "close":
			waiter.Dec()
		case "get-digest":
			if len(req.Args) == 0{
				return Reply{nil, errors.New("expected folder")}
			}
			folder, ok := req.Args[0].(string)
			if !ok {
				return Reply{nil, errors.New("expected folder")}
			}
			digest, err := syncer.GetDigest(folder)
			return Reply{digest, err}
		case "get-file":
			if len(req.Args) < 1 {
				return Reply{nil, errors.New("expected folder")}
			}
			folder, ok := req.Args[0].(string)
			if !ok {
				return Reply{nil, errors.New("expected folder")}
			}

			if len(req.Args) < 2 {
				return Reply{nil, errors.New("expected file")}
			}
			file, ok := req.Args[1].(string)
			if !ok {
				return Reply{nil, errors.New("expected file")}
			}

			stream := make(chan []byte)
			go syncer.ReadFile(folder, file, stream)
			return Reply{stream, nil}
		case "sync":
			if len(req.Args) == 0{
				return Reply{nil, errors.New("expected folder")}
			}
			folder, ok := req.Args[0].(string)
			if !ok {
				return Reply{nil, errors.New("expected folder")}
			}
			waiter.Inc()
			err := dosync(talk, folder)
			return Reply{nil, err}
		}
		return Reply{nil,nil}
	})
}

func handle(conn net.Conn) {
	log.Println("CONNECTION", conn.RemoteAddr())

	waiter := NewWaiter()
	talk := NewTalk(conn)
	waiter.Inc()
	go recv(conn, talk, waiter)
	waiter.Wait()
}

func handleSync(conn net.Conn) {
	waiter := NewWaiter()
	talk := NewTalk(conn)
	go talk.Send(Request{"sync",[]interface{}{"tmp"}})
	waiter.Inc()
	go func() {
		dosync(talk, "tmp")
		waiter.Dec()
	}()
	waiter.Inc()
	go recv(conn, talk, waiter)
	waiter.Wait()
}

func main() {
	flag.Parse()

	from := flag.Arg(0)
	to := flag.Arg(1)

	syncer.Folders["tmp"] = "D:\\tmp\\" + strings.Split(from, ":")[1] + "\\"

	go listen(from)

	r := bufio.NewReader(os.Stdin)
	for {
		ln, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fs := strings.Fields(ln)
		if len(fs) == 0 {
			continue
		}
		switch fs[0] {
		case "sync":
			conn, err := net.Dial("tcp", to)
			if err != nil {
				log.Fatal(err)
			}
			handleSync(conn)
		}
	}
}