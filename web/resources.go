package web

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type (
	ReadSeekerCloser interface {
		io.Reader
		io.Seeker
		io.Closer
	}
	byteSeekerCloser struct {
		*bytes.Reader
	}
	Resource interface {
		Open() (ReadSeekerCloser, error)
		Hash() string
		Name() string
	}
	BaseResource struct {
		path string
		hash string
	}
	FileResource struct {
		BaseResource
	}
	InMemoryResource struct {
		BaseResource
		contents []byte
	}
)

func (this byteSeekerCloser) Close() error {
	return nil
}

func init() {
	root := "/~/"
	http.HandleFunc(root, func(w http.ResponseWriter, req *http.Request) {
		ctx := &defaultContext{DefaultApplication, w, req, "", ""}
			
		components := strings.SplitN(req.URL.Path[len(root):], "/", 2)
		// Make sure we have what we need
		if len(components) != 2 {
			DefaultApplication.ErrorHandler(ctx, errors.New("hash and path are required"))
			return
		}
		h := components[0]
		p := components[1]
		m := getIndex()
		if r, ok := m[h]; ok {
			rsc, err := r.Open()
			if err != nil {
				DefaultApplication.ErrorHandler(ctx, errors.New("hash and path are required"))
				return
			}
			defer rsc.Close()
			
			ctx.Response().Header().Set("Expires", time.Now().Add(365 * 24 * time.Hour).UTC().Format(time.RFC1123))			
			ctx.Response().Header().Set("Cache-Control", "public, max-age=3153600")
			
			http.ServeContent(w, req, p, time.Now().Add(-365 * 24 * time.Hour), rsc)
		} else {			
			DefaultApplication.MissingHandler(ctx)
		}
	})
}

func (this *BaseResource) Hash() string {
	return this.hash
}
func (this *BaseResource) Name() string {
	return this.path
}

func (this *FileResource) Open() (ReadSeekerCloser, error) {
	return os.Open(this.path)
}

func GetResource(fn string) (Resource, error) {
	reader, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	
	return &FileResource{
		BaseResource{
			fn,
			GetHash(reader),
		},
	}, nil
}

func GetHash(reader io.Reader) string {
	h := fnv.New64()
	io.Copy(h, reader)
	key := hex.EncodeToString(h.Sum(nil))
	return key
}

func getIndex() map[string]Resource {
	m := make(map[string]Resource)
	roots := []string{"app/", "vendor/"}
	for _, root := range roots {
		filepath.Walk(root, func(fn string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			
			if err != nil {
				return err
			}
			
			switch path.Ext(fn) {
			case ".go": fallthrough
			case ".tpl":
				return nil
			default:
				r, err := GetResource(fn)
				if err != nil {
					return err
				}
				m[r.Hash()] = r
			}
			
			return nil
		})
	}
	return m
}

func findScript(name string) string {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	dirs := []string{"vendor/js/", "app/js/"}
	exts := []string{".js", ".min.js", "-min.js"}
	
	for _, d := range dirs {
		for _, e := range exts {
			_, err := os.Stat(d + name + e)
			if err == nil {
				return d + name + e
			}
		}
	}
	
	return ""
}
func findStyle(name string) string {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	dirs := []string{"vendor/css/", "app/css/"}
	exts := []string{".css", ".min.css"}
	
	for _, d := range dirs {
		for _, e := range exts {
			_, err := os.Stat(d + name + e)
			if err == nil {
				return d + name + e
			}
		}
	}
	
	return ""
}
func getScriptResources(script string) ([]Resource, error) {
	// Get the resource for this file
	resource, err := GetResource(script)
	if err != nil {
		return nil, err
	}
	// Open it
	handle, err := resource.Open()
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	
	dependencies := []Resource{}
	reader := bufio.NewReader(handle)	
	for {
		// Read the file line by line
		ln, err := reader.ReadString('\n')
		if err != io.EOF && err != nil {
			return nil, err
		}
		ln = strings.TrimSpace(ln)
		// look for comments that start with = require
		if strings.HasPrefix(ln, "//=") || strings.HasPrefix(ln, "*=") {
			fields := strings.Fields(ln)
			if len(fields) == 3 && fields[1] == "require" {
				fn := findScript(fields[2])
				if fn != "" {
					children, err := getScriptResources(fn)
					if err != nil {
						return nil, err
					}
					dependencies = append(dependencies, children...)
				}
			}
		}
		if err == io.EOF {
			break
		}
	}
	return append(dependencies, resource), nil
}

func getScriptURLs(root string) ([]string, error) {
	urls := []string{}
	
	nms := []string{root + "application"}
	for _, nm := range nms {
		fn := findScript(nm)
		if fn != "" {
			resources, err := getScriptResources(fn)
			if err != nil {
				return nil, err
			}
			
			for _, resource := range resources {
				urls = append(urls, "/~/" + resource.Hash() + "/" + resource.Name())
			}
		}
	}
	
	return urls, nil
}
func getStyleResources(style string) ([]Resource, error) {
	// Get the resource for this file
	resource, err := GetResource(style)
	if err != nil {
		return nil, err
	}
	// Open it
	handle, err := resource.Open()
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	
	dependencies := []Resource{}
	reader := bufio.NewReader(handle)	
	for {
		// Read the file line by line
		ln, err := reader.ReadString('\n')
		if err != io.EOF && err != nil {
			return nil, err
		}
		ln = strings.TrimSpace(ln)
		// look for comments that start with = require
		if strings.HasPrefix(ln, "*=") {
			fields := strings.Fields(ln)
			if len(fields) == 3 && fields[1] == "require" {
				fn := findStyle(fields[2])
				if fn != "" {
					children, err := getStyleResources(fn)
					if err != nil {
						return nil, err
					}
					dependencies = append(dependencies, children...)
				}
			}
		}
		if err == io.EOF {
			break
		}
	}
	return append(dependencies, resource), nil
}

func getStyleURLs(root string) ([]string, error) {
	urls := []string{}
	
	nms := []string{root + "application"}
	for _, nm := range nms {
		fn := findStyle(nm)
		if fn != "" {
			resources, err := getStyleResources(fn)
			if err != nil {
				return nil, err
			}
			
			for _, resource := range resources {
				urls = append(urls, "/~/" + resource.Hash() + "/" + resource.Name())
			}
		}
	}
	
	return urls, nil
}