gen-protoc:
	protoc -I=. --go_out=. pkg/hashi/messages.proto
