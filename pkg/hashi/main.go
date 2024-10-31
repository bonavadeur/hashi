package hashi

type BridgeType string

const (
	HASHI_TYPE_SYNC_CLIENT       BridgeType = "sync-client"
	HASHI_TYPE_SYNC_SERVER       BridgeType = "sync-server"
	HASHI_TYPE_HALF_ASYNC_CLIENT BridgeType = "half-async-client"
	HASHI_TYPE_HALF_ASYNC_SERVER BridgeType = "half-async-server"
)

type BridgeCallback func(params ...interface{}) (interface{}, error)
