package router

import (
	"encoding/json"
	"github.com/Postcord/rest"
)

type tapeItem struct {
	// Input
	FuncName string			   `json:"func_name"`
	Params   []json.RawMessage `json:"params"`

	// Output
	Results		 []json.RawMessage `json:"results"`
	GenericError string            `json:"generic_error,omitempty"`
	RESTError    *rest.ErrorREST   `json:"rest_error,omitempty"`
}

type tape []*tapeItem

func (t *tape) write(funcName string, params ...interface{}) *tapeItem {
	p := make([]json.RawMessage, len(params))
	for i, x := range params {
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		p[i] = b
	}
	x := &tapeItem{
		FuncName: funcName,
		Params:   p,
	}
	*t = append(*t, x)
	return x
}

func (t *tapeItem) end(items ...interface{}) {
	// Check if the last type is an error and if so split it from the items.
	var err error
	var ok bool
	if len(items) > 0 {
		err, ok = items[len(items)-1].(error)
		if ok || items[len(items)-1] == nil {
			items = items[:len(items)-1]
		}
	}

	// Marshal the rest of the results.
	p := make([]json.RawMessage, len(items))
	for i, x := range items {
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		p[i] = b
	}
	t.Results = p

	// Figure out how to process the error.
	if err != nil {
		if e, ok := err.(*rest.ErrorREST); ok {
			t.RESTError = e
		} else {
			t.GenericError = err.Error()
		}
	}
}
