package hashi

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
)

type HalfAsyncHashi struct {
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
	sendLock       sync.Mutex
	BucketChan     chan proto.Message
}

func NewHalfAsyncHashi(
	name string,
	bridgeType BridgeType,
	upstreamFile string,
	downstreamFile string,
	requestSchema reflect.Type,
	responseSchema reflect.Type,
	serverCallback BridgeCallback,
) *HalfAsyncHashi {
	newHalfAsyncHashi := &HalfAsyncHashi{
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
		sendLock:       sync.Mutex{},
		BucketChan:     make(chan proto.Message, 1),
	}

	var err error
	checkPipeExist(downstreamFile)
	checkPipeExist(upstreamFile)

	if bridgeType == BRIDGE_TYPE_ASYNC_SERVER {
		newHalfAsyncHashi.downstream, err = os.OpenFile(downstreamFile, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		newHalfAsyncHashi.upstream, err = os.OpenFile(upstreamFile, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		go newHalfAsyncHashi.AsyncReceiveServer()
	}
	if bridgeType == BRIDGE_TYPE_ASYNC_CLIENT {
		newHalfAsyncHashi.upstream, err = os.OpenFile(upstreamFile, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		newHalfAsyncHashi.downstream, err = os.OpenFile(downstreamFile, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}
		go newHalfAsyncHashi.AsyncReceiveClient()
	}

	return newHalfAsyncHashi
}

func (hah *HalfAsyncHashi) AsyncSendClient(message proto.Message) (proto.Message, error) { // for Client
	// marshal message
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatalln("Failed to encode sentMessage:", err)
		return nil, err
	}

	// increse MessageIDCount by 1 and prepare sentMessage
	hah.mu.Lock()
	messageID := hah.increaseMessageIDCount()
	hah.mu.Unlock()

	// prepare sentMessage
	sentMessage := []byte{}
	sentMessage = append(sentMessage, strconv.Itoa(int(messageID))...)
	sentMessage = append(sentMessage, byte(0))
	sentMessage = append(sentMessage, messageBytes...)

	// send
	hah.sendLock.Lock()
	_, err = hah.upstream.Write(sentMessage)
	if err != nil {
		fmt.Println("Error writing to upstream:", err)
		return nil, err
	}

	// receive
	receivedMessage := <-hah.BucketChan
	hah.sendLock.Unlock()

	return receivedMessage, nil
}

func (hah *HalfAsyncHashi) AsyncReceiveClient() { // for Client
	// receive
	for {
		n, err := hah.downstream.Read(hah.buffer)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error reading from FIFO: %v", err)
			panic(err)
		}

		zeroPos := findPositionOfZero(hah.buffer[:n])
		receivedMessageBytes := hah.buffer[zeroPos+1 : n]

		receivedMessage := reflect.New(hah.responseSchema).Interface().(proto.Message)
		err = proto.Unmarshal(receivedMessageBytes, receivedMessage)
		if err != nil {
			log.Fatalf("AsyncReceiveClient Failed to unmarshal message: %v", err)
			panic(err)
		}

		hah.BucketChan <- receivedMessage
	}
}

func (hah *HalfAsyncHashi) AsyncReceiveServer() error { // for Server
	// receive
	for {
		n, err := hah.downstream.Read(hah.buffer)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error reading from FIFO: %v", err)
			return err
		}

		zeroPos := findPositionOfZero(hah.buffer[:n])
		messageID, _ := strconv.Atoi(string(hah.buffer[:zeroPos])) // int
		receivedMessageBytes := hah.buffer[zeroPos+1 : n]

		fmt.Println("received message:", messageID)

		receivedMessage := reflect.New(hah.requestSchema).Interface().(proto.Message)
		err = proto.Unmarshal(receivedMessageBytes, receivedMessage)
		if err != nil {
			log.Fatalf("AsyncReceiveServer Failed to unmarshal message: %v", err)
			return err
		}

		go hah.AsyncSendServer(uint32(messageID))
		hah.BucketChan <- receivedMessage
	}
}

func (hah *HalfAsyncHashi) AsyncSendServer(messageID uint32) { // for Server
	receivedMessage := <-hah.BucketChan

	// run callback function
	result, _ := hah.serverCallback(receivedMessage)

	// marshal message
	resultBytes, err := proto.Marshal(result.(proto.Message))
	if err != nil {
		log.Fatalln("Failed to encode sentMessage:", err)
		panic(err)
	}

	// increse MessageIDCount by 1 and prepare sentMessage
	responseMessage := []byte{}
	responseMessage = append(responseMessage, strconv.Itoa(int(messageID))...)
	responseMessage = append(responseMessage, byte(0))
	responseMessage = append(responseMessage, resultBytes...)

	// send
	hah.sendLock.Lock()
	_, err = hah.upstream.Write(responseMessage)
	fmt.Println("sent message", messageID)
	if err != nil {
		fmt.Println("Error writing to upstream:", err)
		panic(err)
	}
	hah.sendLock.Unlock()
}

func (hah *HalfAsyncHashi) increaseMessageIDCount() uint32 {
	if hah.MessageIDCount == math.MaxUint32 {
		hah.MessageIDCount = 0
	} else {
		hah.MessageIDCount++
	}
	return hah.MessageIDCount
}
