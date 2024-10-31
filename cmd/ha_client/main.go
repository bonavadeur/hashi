package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

func main() {
	// start timing
	start := time.Now()

	client := hashi.NewHalfAsyncHashi(
		"async-client",
		hashi.HASHI_TYPE_HALF_ASYNC_CLIENT,
		"/tmp/client-server",
		"/tmp/server-client",
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Request{}),
		nil,
	)

	sentMessage := &hashi.Request{
		SourceIP: "192.168.101.117",
		Domain:   "hello.default.svc.cluster.local",
		URI:      "/",
		Method:   "GET",
		Headers: []*hashi.Request_Header{
			{
				Field: "ram",
				Value: "100",
			},
			{
				Field: "time",
				Value: "2000",
			},
		},
	}

	for i := 0; i < 10000; i++ {
		go func() {
			_, err := client.AsyncSendClient(sentMessage)
			if err != nil {
				panic(err)
			}
		}()
	}

	// End timing
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
	time.Sleep(3 * time.Second)
}
