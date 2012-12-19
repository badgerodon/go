package main

import (
	"encoding/gob"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"path"
	"path/filepath"
	"os"
	"strings"
)

type (
	DigestEntry struct {
		Path string
		Mode os.FileMode
		Time int64
		Size uint64
		Hash uint64
	}
	Digest struct {
		Entries []DigestEntry
	}

	Syncer struct {
		Folders map[string]string
	}
)

func init() {
	gob.Register(new(DigestEntry))
	gob.Register(new(Digest))
}

func NewSyncer() *Syncer {
	return &Syncer{make(map[string]string)}
}
/*
func (this *Syncer) StartSync(args StartSyncArgs, reply *bool) error {
	client, err := rpc.Dial("tcp", args.Remote)
	if err != nil {
		return err
	}
	defer client.Close()

	var remote *Digest
	err = client.Call("Syncer.GetDigest", args.Folder, &remote)
	if err != nil {
		return err
	}

	var local *Digest
	err = this.GetDigest(args.Folder, &local)
	if err != nil {
		return err
	}

	var ok bool
	err = this.HandleSync(
		HandleSyncArgs{args.Folder, args.Remote, local, remote},
		&ok,
	)
	if err != nil {
		return err
	}

	err = client.Call(
		"Syncer.HandleSync",
		HandleSyncArgs{args.Folder, args.Local, nil, local},
		&ok,
	)
	if err != nil {
		return err
	}

	return nil
}
func (this *Syncer) HandleSync(args HandleSyncArgs, reply *bool) error {
	local := args.Local
	if local == nil {
		err := this.GetDigest(args.Folder, &local)
		if err == nil {
			return err
		}
	}
	remote := args.Remote

	toget := []string{}
	todel := []string{}

	i := 0
	j := 0
	for {
		if !(i < len(local.Entries) || j < len(local.Entries)) {
			break
		}
		// Not in local, get file
		if i >= len(local.Entries) {
			toget = append(toget, remote.Entries[j].Path)
			j++
			continue
		}
		// Not in remote
		if j >= len(local.Entries) {
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

	}
}*/
func (this *Syncer) GetDigest(folder string) (*Digest, error) {
	dir, ok := this.Folders[folder]
	if !ok {
		return nil, errors.New("Unknown folder `" + folder + "`")
	}
	if dir[len(dir)-1] == '\\' || dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}

	entries := make([]DigestEntry, 0)
	err := filepath.Walk(dir, func(fn string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Fix the filename so it's consistent across machines
		rel := strings.Replace(fn, "\\", "/", -1)
		rel = rel[len(strings.Replace(dir, "\\", "/", -1)):]
		if info.IsDir() {
			rel += "/"
			entries = append(entries, DigestEntry{
				Path: rel,
				Mode: info.Mode(),
				Time: info.ModTime().UnixNano(),
				Size: 0,
				Hash: 0,
			})
			return nil
		}

		// Open the file for reading
		reader, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer reader.Close()

		h := fnv.New64()
		io.Copy(h, reader)

		// Append the entry
		entries = append(entries, DigestEntry{
			Path: rel,
			Mode: info.Mode(),
			Time: info.ModTime().UnixNano(),
			Size: uint64(info.Size()),
			Hash: h.Sum64(),
		})
		return nil
	})

	return &Digest{entries}, err
}
func (this *Syncer) ReadFile(folder, file string, stream chan []byte) error {
	defer close(stream)

	dir, ok := this.Folders[folder]
	if !ok {
		return errors.New("Unknown folder `" + folder + "`")
	}
	fn := path.Join(dir, file)

	r, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer r.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		stream <- buf[:n]
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
func (this *Syncer) WriteFile(folder, file string, stream chan []byte) error {
	dir, ok := this.Folders[folder]
	if !ok {
		return errors.New("Unknown folder `" + folder + "`")
	}
	fn := path.Join(dir, file)

	w, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer w.Close()

	for {
		bs, ok := <- stream
		if !ok {
			break
		}
		log.Println("BS", bs)
		_, err = w.Write(bs)
		if err != nil {
			return err
		}
	}

	return nil
}