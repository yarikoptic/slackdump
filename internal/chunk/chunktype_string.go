// Code generated by "stringer -type=ChunkType -trimprefix=C"; DO NOT EDIT.

package chunk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CMessages-0]
	_ = x[CThreadMessages-1]
	_ = x[CFiles-2]
}

const _ChunkType_name = "MessagesThreadMessagesFiles"

var _ChunkType_index = [...]uint8{0, 8, 22, 27}

func (i ChunkType) String() string {
	if i < 0 || i >= ChunkType(len(_ChunkType_index)-1) {
		return "ChunkType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChunkType_name[_ChunkType_index[i]:_ChunkType_index[i+1]]
}
