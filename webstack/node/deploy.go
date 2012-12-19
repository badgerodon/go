package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"net/rpc"
)

func deploy(args []string) {
	if len(args) == 0 {
		fmt.Println("Expected server name")
		os.Exit(1)
	}

	conn, err := tls.Dial("tcp", args[0], &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte(PASSWORD))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	
	pkg, err := ReadPackage(".")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	crwc := NewCompressedReadWriteCloser(conn)
	client := rpc.NewClient(crwc)

	fmt.Println("Deploying...")
	var status bool
	err = client.Call("Deployer.Deploy", &pkg, &status)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(status)
}