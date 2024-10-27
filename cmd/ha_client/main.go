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
		hashi.BRIDGE_TYPE_ASYNC_CLIENT,
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

	for i := 0; i < 100; i++ {
		go func() {
			result, err := client.AsyncSendClient(sentMessage)
			if err != nil {
				panic(err)
			}
			fmt.Println("result", result, "\n")
		}()
	}

	// End timing
	elapsed := time.Since(start)
	// fmt.Println("received:", receivedMessage)
	fmt.Printf("Time taken: %s\n", elapsed)
	time.Sleep(3 * time.Second)
}
