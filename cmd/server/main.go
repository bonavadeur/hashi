package main

import (
	"context"
	"reflect"
	"time"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

type Nani struct {
	messages []interface{}
}

func (n *Nani) Callback(params ...interface{}) (interface{}, error) {
	time.Sleep(5 * time.Millisecond)
	n.messages = append(n.messages, params[0])
	return params[0], nil
}

func main() {
	ctx := context.Background()

	nani := &Nani{
		messages: []interface{}{},
	}

	_ = hashi.NewHashi(
		"server",
		hashi.HASHI_TYPE_SERVER,
		"/tmp/hashi",
		10,
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Request{}),
		nani.Callback,
	)

	<-ctx.Done()
}
