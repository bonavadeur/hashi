package main

import (
	"context"
	"reflect"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

func main() {
	ctx := context.Background()

	callback := func(params ...interface{}) (interface{}, error) {
		_ = params[0]
		// fmt.Println(message)
		return nil, nil
	}

	_ = hashi.NewSyncHashi(
		"server",
		hashi.HASHI_TYPE_SYNC_SERVER,
		"/tmp/server-client",
		"/tmp/client-server",
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Response{}),
		callback,
	)

	<-ctx.Done()
}
