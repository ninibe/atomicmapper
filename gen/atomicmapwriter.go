package gen

import (
	"io"

	"github.com/clipperhouse/typewriter"
)

func init() {
	err := typewriter.Register(NewAtomicMapWriter())
	if err != nil {
		panic(err)
	}
}

type AtomicMapWriter struct{}

func NewAtomicMapWriter() *AtomicMapWriter {
	return &AtomicMapWriter{}
}

func (sw *AtomicMapWriter) Name() string {
	return "atomicmap"
}

func (sw *AtomicMapWriter) Imports(t typewriter.Type) (result []typewriter.ImportSpec) {
	return []typewriter.ImportSpec{
		{Path: "sync"},
		{Path: "sync/atomic"},
	}
}

func (sw *AtomicMapWriter) Write(w io.Writer, t typewriter.Type) error {
	tmpl, err := template.Parse()
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, t); err != nil {
		return err
	}

	return nil
}

var template = &typewriter.Template{
	Name: "AtomicMap",
	Text: MapTPL,
}

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
func (am *{{.Name}}AtomicMap) Get(key string) {{.Pointer}}{{.Name}} {
	return am.val.Load().(_{{.Name}}Map)[key]
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
