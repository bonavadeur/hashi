package hashi

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Hashi struct {
	Name           string
	bridgeType     BridgeType // ["client", "server"]
	pipeFolder     string
	concurrency    int
	requestSchema  reflect.Type
	responseSchema reflect.Type
	serverCallback BridgeCallback
	chisanaHashis  []*HalfAsyncHashi
	MessageIDCount int
	mIDLock        sync.Mutex
	bucket         chan proto.Message
}

func NewHashi(
	name string,
	bridgeType BridgeType,
	pipeFolder string,
	concurrency int,
	requestSchema reflect.Type,
	responseSchema reflect.Type,
	serverCallback BridgeCallback,
) *Hashi {
	newHashi := &Hashi{
		Name:           name,
		bridgeType:     bridgeType,
		pipeFolder:     pipeFolder,
		concurrency:    concurrency,
		requestSchema:  requestSchema,
		responseSchema: responseSchema,
		serverCallback: serverCallback,
		chisanaHashis:  []*HalfAsyncHashi{},
		MessageIDCount: 0,
		mIDLock:        sync.Mutex{},
		bucket:         make(chan proto.Message, 1),
	}

	for i := 0; i < concurrency; i++ {
		var _upstreamFile, _downstreamFile string
		var _bridgeType BridgeType
		if bridgeType == HASHI_TYPE_CLIENT {
			_upstreamFile = fmt.Sprintf("%s/%d/%s", pipeFolder, i, "client-server")
			_downstreamFile = fmt.Sprintf("%s/%d/%s", pipeFolder, i, "server-client")
			_bridgeType = HASHI_TYPE_HALF_ASYNC_CLIENT
		}
		if bridgeType == HASHI_TYPE_SERVER {
			_upstreamFile = fmt.Sprintf("%s/%d/%s", pipeFolder, i, "server-client")
			_downstreamFile = fmt.Sprintf("%s/%d/%s", pipeFolder, i, "client-server")
			_bridgeType = HASHI_TYPE_HALF_ASYNC_SERVER
		}

		newHashi.chisanaHashis = append(
			newHashi.chisanaHashis,
			NewHalfAsyncHashi(
				strconv.Itoa(i),
				_bridgeType,
				_upstreamFile,
				_downstreamFile,
				requestSchema,
				responseSchema,
				serverCallback,
			),
		)
	}

	return newHashi
}

func (h *Hashi) AsyncSendClient(message proto.Message) (proto.Message, error) { // for Client
	// increse MessageIDCount by 1 and prepare sentMessage
	h.mIDLock.Lock()
	messageID := h.increaseMessageIDCount()
	h.mIDLock.Unlock()

	// send and receive
	result, err := h.chisanaHashis[messageID-1].AsyncSendClient(message)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Hashi) increaseMessageIDCount() int {
	if h.MessageIDCount == h.concurrency {
		h.MessageIDCount = 1
	} else {
		h.MessageIDCount++
	}
	return h.MessageIDCount
}
