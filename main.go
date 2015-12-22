package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var name = flag.String("type", "", "[required] type to store in the map")
var point = flag.Bool("pointer", false, "store and retrieve pointers to the type")
var out = flag.String("out", "", "output file name")
var pack = flag.String("package", "", "go package of the new file")

func main() {
	log.SetFlags(0)
	log.SetPrefix("atomicmapper: ")
	flag.Parse()

	if *name == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	tmpl, err := template.New("mapper").Parse(MapTPL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	outName := *out
	if outName == "" {
		outName = fmt.Sprintf("%s_atomicmap.go", strings.ToLower(*name))
	}

	packName := os.Getenv("GOPACKAGE")
	if *pack != "" {
		packName = *pack
	}

	if packName == "" {
		here, _ := filepath.Abs(flag.Args()[0])
		packName = filepath.Base(filepath.Dir(here))
	}

	f, err := os.Create(outName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	f.WriteString("// generated file - DO NOT EDIT\n\n\n")
	f.WriteString("package " + packName + "\n")
	f.WriteString(ImportTPL)

	var Pointer string
	if *point {
		Pointer = "*"
	}

	err = tmpl.Execute(f, struct {
		Name string
		Pointer string
	}{
		Name: *name,
		Pointer: Pointer,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

const ImportTPL = `
import (
	"sync"
	"sync/atomic"
)
`

const MapTPL = `
// {{.Name}}AtomicMap is a copy-on-write thread-safe map of {{if .Pointer}}pointers to {{end}}{{.Name}}
type {{.Name}}AtomicMap struct {
	mu sync.Mutex
	val atomic.Value
}

type _{{.Name}}Map map[string]{{.Pointer}}{{.Name}}

// New{{.Name}}AtomicMap returns a new initialized {{.Name}}AtomicMap
func New{{.Name}}AtomicMap() *{{.Name}}AtomicMap {
	am := &{{.Name}}AtomicMap{}
	am.val.Store(make(_{{.Name}}Map, 0))
	return am
}

// Get returns a {{if .Pointer}}pointer to {{end}}{{.Name}} for a given key
func (am *{{.Name}}AtomicMap) Get(key string) (value {{.Pointer}}{{.Name}}, ok bool) {
	value, ok = am.val.Load().(_{{.Name}}Map)[key]
	return value, ok
}

// Len returns the number of elements in the map
func (am *TypeAtomicMap) Len() int {
	return len(am.val.Load().(_{{.Name}}Map))
}

// Set inserts in the map a {{if .Pointer}}pointer to {{end}}{{.Name}} under a given key
func (am *{{.Name}}AtomicMap) Set(key string, value {{.Pointer}}{{.Name}}) {
	am.mu.Lock()
	defer am.mu.Unlock()

	m1 := am.val.Load().(_{{.Name}}Map)
	m2 := make(_{{.Name}}Map, len(m1)+1)
	for k, v := range m1 {
		m2[k] = v
	}

	m2[key] = value
	am.val.Store(m2)
	return
}

// Delete removes the {{if .Pointer}}pointer to {{end}}{{.Name}} under key from the map
func (am *{{.Name}}AtomicMap) Delete(key string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	m1 := am.val.Load().(_{{.Name}}Map)
	_, ok := m1[key]
	if !ok {
		return
	}

	m2 := make(_{{.Name}}Map, len(m1)-1)
	for k, v := range m1 {
		if k != key {
			m2[k] = v
		}
	}

	am.val.Store(m2)
	return
}
`
