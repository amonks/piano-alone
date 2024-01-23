// Code generated by "stringer -type=MessageType"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MessageTypeInvalid-0]
	_ = x[MessageTypeJoin-1]
	_ = x[MessageTypeLeave-2]
	_ = x[MessageTypeSubmitPartialTrack-3]
	_ = x[MessageTypeInitialState-4]
	_ = x[MessageTypeBroadcastConnectedPlayer-5]
	_ = x[MessageTypeBroadcastDisconnectedPlayer-6]
	_ = x[MessageTypeAssignment-7]
	_ = x[MessageTypeBroadcastPhase-8]
	_ = x[MessageTypeBroadcastCombinedTrack-9]
}

const _MessageType_name = "MessageTypeInvalidMessageTypeJoinMessageTypeLeaveMessageTypeSubmitPartialTrackMessageTypeInitialStateMessageTypeBroadcastConnectedPlayerMessageTypeBroadcastDisconnectedPlayerMessageTypeAssignmentMessageTypeBroadcastPhaseMessageTypeBroadcastCombinedTrack"

var _MessageType_index = [...]uint8{0, 18, 33, 49, 78, 101, 136, 174, 195, 220, 253}

func (i MessageType) String() string {
	if i >= MessageType(len(_MessageType_index)-1) {
		return "MessageType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MessageType_name[_MessageType_index[i]:_MessageType_index[i+1]]
}
