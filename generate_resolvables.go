//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"
)

const start = `// Code generated by generate_resolvables.go; DO NOT EDIT.

package router

//go:generate go run generate_resolvables.go

import (
	"encoding/json"
	"strconv"

	"github.com/Postcord/objects"
)

`

const singleStructureTemplate = `// Resolvable{{ .Type }} is used to define a {{ .Type }} in a command option that is potentially resolvable.
type Resolvable{{ .Type }} struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

// Snowflake is used to return the ID as a snowflake.
func (r Resolvable{{ .Type }}) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

// MarshalJSON implements the json.Marshaler interface.
func (r Resolvable{{ .Type }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// String is used to return the ID as a string.
func (r Resolvable{{ .Type }}) String() string {
	return r.id
}

// Resolve is used to attempt to resolve the item to its type. Returns nil if it doesn't exist.
func (r Resolvable{{ .Type }}) Resolve() *objects.{{ .Type }} {
	x, ok := r.data.Resolved.{{ .Type }}s[r.Snowflake()]
	if !ok {
		return nil
	}
	return &x
}`

var types = []string{
	"User", "Channel", "Role", "Message", "Attachment",
}

func main() {
	file := start
	parts := make([]string, len(types))
	t, err := template.New("_").Parse(singleStructureTemplate)
	if err != nil {
		panic(err)
	}
	for i, v := range types {
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, map[string]interface{}{"Type": v}); err != nil {
			panic(err)
		}
		parts[i] = buf.String()
	}
	file += strings.Join(parts, "\n\n") + "\n"
	if err := ioutil.WriteFile("resolvables_gen.go", []byte(file), 0666); err != nil {
		panic(err)
	}
}
