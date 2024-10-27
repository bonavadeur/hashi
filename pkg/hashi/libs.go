package hashi

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"

	"google.golang.org/protobuf/proto"
)

func checkPipeExist(pipePath string) {
	if _, err := os.Stat(pipePath); os.IsNotExist(err) {
		if err := syscall.Mkfifo(pipePath, 0777); err != nil {
			fmt.Println("Error creating named pipe:", err)
			return
		}
	}
}

func setField(message proto.Message, fieldName string, value interface{}) {
	v := reflect.ValueOf(message).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		log.Fatalf("No such field: %s in message", fieldName)
	}
	if !field.CanSet() {
		log.Fatalf("Cannot set field: %s in message", fieldName)
	}
	field.Set(reflect.ValueOf(value))
}

func findPositionOfZero(data []byte) int {
	for i, b := range data {
		if b == byte(0) {
			return i
		}
	}
	return -1
}
