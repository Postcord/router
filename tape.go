package router

import "encoding/json"

type tapeItem struct {
	FuncName string			   `json:"func_name"`
	Params   []json.RawMessage `json:"params"`
}

type tape []tapeItem

func (t *tape) write(funcName string, params ...interface{}) {
	p := make([]json.RawMessage, len(params))
	for i, x := range params {
		b, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		p[i] = b
	}
	*t = append(*t, tapeItem{
		FuncName: funcName,
		Params:   p,
	})
}
