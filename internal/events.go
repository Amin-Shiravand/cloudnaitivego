package src

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)
