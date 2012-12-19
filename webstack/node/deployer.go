package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type (
	FileInfo struct {
		Size int64
		Mode os.FileMode
		ModTime time.Time
	}
	PackageEntry struct {
		Path string
		Info FileInfo
		Contents []byte
	}
	Package []PackageEntry
	Deployer struct {}
)

func init() {
	gob.Register(FileInfo{})
	gob.Register(PackageEntry{})
	gob.Register(Package{})
}

func ReadPackage(dir string) (Package, error) {
	pkg := []PackageEntry{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		contents, e := ioutil.ReadFile(p)
		if e != nil {
			return e
		}
		pkg = append(pkg, PackageEntry{
			strings.Replace(p, "\\", "/", -1),
			FileInfo{info.Size(),info.Mode(),info.ModTime()},
			contents,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return Package(pkg), nil
}
func WritePackage(dir string, pkg Package) error {
	for _, entry := range *pkg {
		ep := path.Join(p, entry.Path)
		ed := path.Dir(ep)
		err = os.MkdirAll(ed, os.ModeDir)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(ep, entry.Contents, entry.Info.Mode)
		if err != nil {
			return err
		}
	}
	return nil
}
func ProcessDeploy(dir string) error {
	bs, err := ioutil.ReadFile(path.Join(dir, "config.json"))
	if err != nil {
		return err
	}
	cfg := map[string]string{}
	err = json.Unmarshal(&cfg)
	if err != nil {
		return err
	}
	role, ok := cfg["role"]
	if !ok {
		return errors.New("Expected `role` in config")
	}

	

}

func (this *Deployer) Deploy(pkg *Package, ok *bool) error {
	nm := fmt.Sprint("deployed-", time.Now().UnixNano())
	p := path.Join(os.TempDir(), nm)
	err := os.Mkdir(p, os.ModeDir)
	if err != nil {
		return err
	}
	err = WritePackage(p, *pkg)
	if err != nil {
		return err
	}
	*ok = true
	return nil
}