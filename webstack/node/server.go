package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"net/rpc"
)

func server(args []string) {
	rpc.Register(new(Deployer))

	cert, err := tls.LoadX509KeyPair(CERTIFICATE, KEY)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	l, err := tls.Listen("tcp", fmt.Sprint(":", PORT), &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error", err)
			continue
		}

		sentPassword := make([]byte, len(PASSWORD))
		_, err = conn.Read(sentPassword)
		if err != nil {
			log.Println("Error", err)
			conn.Close()
			continue
		}

		if string(sentPassword) != PASSWORD {
			log.Println("Invalid password")
			conn.Close()
			continue
		}

		crwc := NewCompressedReadWriteCloser(conn)
		go rpc.ServeConn(crwc)
	}
}