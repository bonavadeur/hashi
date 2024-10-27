package hashi

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

type SyncHashi struct {
	Name           string
	bridgeType     BridgeType // ["client", "server"]
	upstreamFile   string
	downstreamFile string
	upstream       *os.File // Write
	downstream     *os.File // Read
	buffer         []byte
	requestSchema  reflect.Type
	responseSchema reflect.Type
	serverCallback BridgeCallback
	MessageIDCount uint32
	mu             sync.Mutex
	
}

func NewSyncHashi(
	name string,
	bridgeType BridgeType,
	upstreamFile string,
	downstreamFile string,
	requestSchema reflect.Type,
	responseSchema reflect.Type,
	serverCallback BridgeCallback,
) *SyncHashi {
	newHashi := &SyncHashi{
		Name:           name,
		bridgeType:     bridgeType,
		upstreamFile:   upstreamFile,
		downstreamFile: downstreamFile,
		MessageIDCount: 0,
		buffer:         make([]byte, 1024),
		requestSchema:  requestSchema,
		responseSchema: responseSchema,
		serverCallback: serverCallback,
		mu:             sync.Mutex{},
	}

	var err error
	checkPipeExist(downstreamFile)
	checkPipeExist(upstreamFile)

	if bridgeType == BRIDGE_TYPE_SYNC_SERVER {
		newHashi.downstream, err = os.OpenFile(downstreamFile, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		newHashi.upstream, err = os.OpenFile(upstreamFile, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		go newHashi.ReceiveAndSend()
	}
	if bridgeType == BRIDGE_TYPE_SYNC_CLIENT {
		newHashi.upstream, err = os.OpenFile(upstreamFile, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		newHashi.downstream, err = os.OpenFile(downstreamFile, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
	}

	return newHashi
}

func (sh *SyncHashi) SendAndReceive(message proto.Message) (proto.Message, error) { // for Client
	// send
	sentMessage, err := proto.Marshal(message)
	if err != nil {
		log.Fatalln("Failed to encode sentMessage:", err)
		return nil, err
	}
	_, err = sh.upstream.Write(sentMessage)
	if err != nil {
		fmt.Println("Error writing to upstream:", err)
		return nil, err
	}

	// receive
	for {
		n, err := sh.downstream.Read(sh.buffer)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error reading from FIFO: %v", err)
			return nil, err
		}

		// receivedMessage := &pb.Response{}
		receivedMessage := reflect.New(sh.responseSchema).Interface().(proto.Message)
		err = proto.Unmarshal(sh.buffer[:n], receivedMessage)
		if err != nil {
			log.Fatalf("Failed to unmarshal message: %v", err)
			return nil, err
		}
		return receivedMessage, nil
	}
}

func (sh *SyncHashi) ReceiveAndSend() { // for Server
	// receive
	for {
		n, err := sh.downstream.Read(sh.buffer)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error reading from FIFO: %v", err)
		}

		// receivedMessage := &pb.Message{}
		receivedMessage := reflect.New(sh.requestSchema).Interface().(proto.Message)
		err = proto.Unmarshal(sh.buffer[:n], receivedMessage)
		if err != nil {
			log.Fatalf("Failed to unmarshal message: %v", err)
		}

		// run callback function
		sh.serverCallback(receivedMessage)

		// send
		sentMessage := reflect.New(sh.responseSchema).Interface().(proto.Message)
		setField(sentMessage, "Status", ResponseStatus_Success)
		retMessageByte, err := proto.Marshal(sentMessage)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		}
		_, err = sh.upstream.Write(retMessageByte)
		if err != nil {
			fmt.Println("Error writing to pipeSC:", err)
			return
		}
	}
}

func (sh *SyncHashi) increaseMessageIDCount() uint32 {
	if sh.MessageIDCount == math.MaxUint32 {
		sh.MessageIDCount = 0
	} else {
		sh.MessageIDCount++
	}
	return sh.MessageIDCount
}
