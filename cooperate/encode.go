package main

import (
	"encoding/binary"
	"io"
	"os"
)

func readFile(r io.Reader, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	sz, err := readInt64(r)
	if err != nil {
		return err
	}

	_, err = io.CopyN(w, r, sz)
	return err
}
func readInt64(r io.Reader) (int64, error) {
	i, err := readUint64(r)
	return int64(i), err
}
func readString(r io.Reader) (string, error) {
	sz, err := readUint32(r)
	if err != nil {
		return "", err
	}
	bs := make([]byte, sz)
	_, err = io.ReadFull(r, bs)
	return string(bs), err
}
func readUint32(r io.Reader) (uint32, error) {
	bs := make([]byte, 4)
	_, err := io.ReadFull(r, bs)
	return binary.BigEndian.Uint32(bs), err
}
func readUint64(r io.Reader) (uint64, error) {
	bs := make([]byte, 8)
	_, err := io.ReadFull(r, bs)
	return binary.BigEndian.Uint64(bs), err
}
func writeFile(w io.Writer, path string) error {
	// Get info (useful for size)
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Open the file for reading
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	// Write the length
	err = writeInt64(w, info.Size())
	if err != nil {
		return err
	}

	// Copy the file
	_, err = io.CopyN(w, r, info.Size())
	return err
}
func writeInt64(w io.Writer, v int64) error {
	return writeUint64(w, uint64(v))
}
func writeString(w io.Writer, v string) error {
	err := writeUint32(w, uint32(len(v)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(v))
	return err
}
func writeUint32(w io.Writer, v uint32) error {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	_, err := w.Write(bs)
	return err
}
func writeUint64(w io.Writer, v uint64) error {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	_, err := w.Write(bs)
	return err
}