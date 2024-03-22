package game

//go:generate go run golang.org/x/tools/cmd/stringer -type=ConnectionState
type ConnectionState byte

const (
	ConnectionStateUninitialized ConnectionState = iota
	ConnectionStateDisconnected
	ConnectionStateConnected
)

