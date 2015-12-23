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

var (
	name  = flag.String("type", "", "[required] type to store in the map")
	point = flag.Bool("pointer", false, "store and retrieve pointers to the type")
	out   = flag.String("out", "", "output file name")
	pack  = flag.String("package", "", "go package of the new file")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("atomicmapper: ")
	flag.Parse()

	if *name == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var Imports []string = []string{"sync", "sync/atomic"}

	var Name string
	var Subpackage string

	nameParts := strings.Split(*name, ".")
	Name = nameParts[len(nameParts)-1]
	if len(nameParts) > 1 {
		FullPackage := strings.Join(nameParts[:len(nameParts)-1], ".")
		Imports = append(Imports, FullPackage)
		packParts := strings.Split(FullPackage, "/")
		Subpackage = packParts[len(packParts)-1]
	}

	outName := *out
	if outName == "" {
		outName = fmt.Sprintf("%s_atomicmap.go", strings.ToLower(Name))
	}

	packName := os.Getenv("GOPACKAGE")
	if *pack != "" {
		packName = *pack
	}

	if packName == "" {
		here, _ := filepath.Abs(os.Args[0])
		packName = filepath.Base(filepath.Dir(here))
	}

	f, err := os.Create(outName)
	fatalIf(err)
	f.WriteString("// generated file - DO NOT EDIT\n")
	f.WriteString("// command: " + strings.Join(os.Args, " ") + "\n\n\n")
	f.WriteString("package " + packName + "\n")

	tmplImp, err := template.New("imports").Parse(ExtImportTPL)
	fatalIf(err)
	err = tmplImp.Execute(f, Imports)
	fatalIf(err)

	tmpl, err := template.New("mapper").Parse(MapTPL)
	fatalIf(err)

	var Pointer string
	if *point {
		Pointer = "*"
	}

	if Subpackage != "" {
		Subpackage += "."
	}

	err = tmpl.Execute(f, struct {
		Name       string
		Pointer    string
		Subpackage string
	}{
		Name:       Name,
		Pointer:    Pointer,
		Subpackage: Subpackage,
	})
	fatalIf(err)
}

func fatalIf(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

const MinImportTPL = `
import (
	"sync"
	"sync/atomic"
)
`

const ExtImportTPL = `
import ({{ range $key, $value := . }}
	"{{ $value }}"{{ end }}
)
`

const MapTPL = `
// {{.Name}}AtomicMap is a copy-on-write thread-safe map of {{if .Pointer}}pointers to {{end}}{{.Name}}
type {{.Name}}AtomicMap struct {
	mu sync.Mutex
	val atomic.Value
}

type _{{.Name}}Map map[string]{{.Pointer}}{{.Subpackage}}{{.Name}}

// New{{.Name}}AtomicMap returns a new initialized {{.Name}}AtomicMap
func New{{.Name}}AtomicMap() *{{.Name}}AtomicMap {
	am := &{{.Name}}AtomicMap{}
	am.val.Store(make(_{{.Name}}Map, 0))
	return am
}

// Get returns a {{if .Pointer}}pointer to {{end}}{{.Name}} for a given key
func (am *{{.Name}}AtomicMap) Get(key string) (value {{.Pointer}}{{.Subpackage}}{{.Name}}, ok bool) {
	value, ok = am.val.Load().(_{{.Name}}Map)[key]
	return value, ok
}

// Len returns the number of elements in the map
func (am *{{.Name}}AtomicMap) Len() int {
	return len(am.val.Load().(_{{.Name}}Map))
}

// Set inserts in the map a {{if .Pointer}}pointer to {{end}}{{.Name}} under a given key
func (am *{{.Name}}AtomicMap) Set(key string, value {{.Pointer}}{{.Subpackage}}{{.Name}}) {
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
