package main

import (
	"context"
	"reflect"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

func main() {
	ctx := context.Background()

	callback := func(params ...interface{}) (interface{}, error) {
		message := params[0]
		// fmt.Println(message)
		return message, nil
	}

	_ = hashi.NewHalfAsyncHashi(
		"async-server",
		hashi.BRIDGE_TYPE_ASYNC_SERVER,
		"/tmp/server-client",
		"/tmp/client-server",
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Request{}),
		callback,
	)

	<-ctx.Done()
}
