// Code generated by "stringer -type=ConnectionState"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ConnectionStateUninitialized-0]
	_ = x[ConnectionStateDisconnected-1]
	_ = x[ConnectionStateConnected-2]
}

const _ConnectionState_name = "ConnectionStateUninitializedConnectionStateDisconnectedConnectionStateConnected"

var _ConnectionState_index = [...]uint8{0, 28, 55, 79}

func (i ConnectionState) String() string {
	if i >= ConnectionState(len(_ConnectionState_index)-1) {
		return "ConnectionState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ConnectionState_name[_ConnectionState_index[i]:_ConnectionState_index[i+1]]
}