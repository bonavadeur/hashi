package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/bonavadeur/hashi/pkg/hashi"
)

func main() {
	// start timing
	start := time.Now()

	client := hashi.NewHashi(
		"client",
		hashi.HASHI_TYPE_CLIENT,
		"/tmp/hashi",
		10,
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

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			_, err := client.AsyncSendClient(sentMessage)
			if err != nil {
				panic(err)
			}
			// fmt.Println(result, "\n")
			wg.Done()
		}()
	}
	wg.Wait()

	// End timing
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed/1000)
}
