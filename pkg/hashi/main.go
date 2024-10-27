package hashi

type BridgeType string

const (
	BRIDGE_TYPE_SYNC_CLIENT  BridgeType = "sync-client"
	BRIDGE_TYPE_SYNC_SERVER  BridgeType = "sync-server"
	BRIDGE_TYPE_ASYNC_CLIENT BridgeType = "async-client"
	BRIDGE_TYPE_ASYNC_SERVER BridgeType = "async-server"
)

type BridgeCallback func(params ...interface{}) (interface{}, error)
