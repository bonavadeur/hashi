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

	client := hashi.NewSyncHashi(
		"client",
		hashi.HASHI_TYPE_SYNC_CLIENT,
		"/tmp/client-server",
		"/tmp/server-client",
		reflect.TypeOf(hashi.Request{}),
		reflect.TypeOf(hashi.Response{}),
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

	for i := 0; i < 1000; i++ {
		_, err := client.SendAndReceive(sentMessage)
		if err != nil {
			panic(err)
		}
	}

	// End timing
	elapsed := time.Since(start)
	// fmt.Println("received:", receivedMessage)
	fmt.Printf("Time taken: %s\n", elapsed)
}
