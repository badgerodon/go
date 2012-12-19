package main

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	//"net/http"
	//"net/http/httputil"
	//"path"
	//"sync"

	//"github.com/howeyc/fsnotify"
)
/*
func (this *Main) ListenForChanges(filename string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	this.watcher = w

	go func() {
		for {
			select {
			case ev := <-this.watcher.Event:
				if ev.Name == filename {
					this.LoadConfig(filename)
				}
			}
		}
	}()

	err = this.watcher.Watch(path.Dir(filename))
	return err
}*/

func main() {
	backend := NewBackEnd()
	err := backend.LoadFile("backend.json")
	if err != nil {
		log.Fatalln(err)
	}

	frontend := NewFrontEnd(backend)
	err = frontend.LoadFile("frontend.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = frontend.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
