package main

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"sync"
)

type (
	BackEndConfigEntry struct {
		Domain string `json:"domain"`
		Host string `json:"host"`
		Port int `json:"port"`
	}
	BackEndConfig []BackEndConfigEntry

	BackEnd struct {
		lock sync.Mutex
		next int
		entries map[string][]BackEndConfigEntry
	}
)

func NewBackEnd() *BackEnd {
	return &BackEnd{}
}

func (this *BackEnd) LoadFile(filename string) error {
	log.Println("Loading backend config from", filename)
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var cfg BackEndConfig
	err = json.Unmarshal(contents, &cfg)
	if err != nil {
		return err
	}
	this.Load(cfg)
	return nil
}

func (this *BackEnd) Load(config BackEndConfig) {
	this.lock.Lock()
	defer this.lock.Unlock()

	entries := make(map[string][]BackEndConfigEntry)
	for _, entry := range config {
		es, ok := entries[entry.Domain]
		if !ok {
			es = make([]BackEndConfigEntry, 1)
		}
		entries[entry.Domain] = append(es, entry)
	}
	this.entries = entries
}

func (this *BackEnd) Get(domain string) (BackEndConfigEntry, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	entries, ok := this.entries[domain]
	if ok {
		this.next += 1
		return entries[this.next % len(entries)], true
	}

	return BackEndConfigEntry{}, false
}