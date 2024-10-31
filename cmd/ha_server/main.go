package main

import (
	"context"
	"reflect"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

type Nani struct {
	messages []interface{}
}

func (n *Nani) Callback(params ...interface{}) (interface{}, error) {
	n.messages = append(n.messages, params[0])
	return params[0], nil
}

func main() {
	ctx := context.Background()

	nani := &Nani{
		messages: []interface{}{},
	}

	_ = hashi.NewHalfAsyncHashi(
		"async-server",
		hashi.HASHI_TYPE_HALF_ASYNC_SERVER,
		"/tmp/server-client",
		"/tmp/client-server",
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Request{}),
		nani.Callback,
	)

	<-ctx.Done()
}
