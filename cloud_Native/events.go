package cloud_Native

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)
