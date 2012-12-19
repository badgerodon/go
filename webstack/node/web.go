package main

import (
	"errors"
	"path"
	"sync"
)

type (
	WebServer struct {
		Domain string
		Folder string
		Command string
		Process *os.Process
	}
)

var (
	webServers map[string][]string
	webLock sync.Mutex
)

func GetCommand(folder string, cfg map[string]string) (*exec.Cmd, error) {
	command, ok := cfg["command"]
	if !ok {
		return nil, errors.New("Expected `command`")
	}

	args := strings.Split(command, " ")
	if len(args) == 0 || args[0] == "" {
		return nil, errors.New("Invalid command")
	}

	exe := args[0]
	args = args[1:]

	if strings.StartsWith(exe, "./") {
		exe = path.Join(folder, exe)
	} else {
		exe, err = exec.LookPath(exe)
		if err != nil {
			return nil, err
		}
	}

	return &Cmd{
		Path: exe,
		Args: args,
		Env: fmt.Sprint("PORT=", nextPort()),
		Dir: folder,
	}
}

func ProcessWeb(folder string, cfg map[string]string) error {
	domain, ok := cfg["domain"]
	if !ok {
		return errors.New("Expected `domain`")
	}

	cmd, err := GetCommand(folder, cfg)
	if err != nil {
		return err
	}

	webLock.Lock()
	defer webLock.Unlock()

	err = cmd.Start()
	if err != nil {
		return err
	}

	newServers := []WebServer{
		WebServer{
			Domain: domain,
			Folder: folder,
			Command: command,
			Process: cmd.Process,
		}
	}

	oldServers, ok := webServers[domain]
	if ok {
		for _, server := range oldServers {
			server.Process.Kill()
			os.RemoveAll(server.Folder)
		}
	}
}