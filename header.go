package fbx

import (
	"encoding/binary"
	"fmt"
)

type Header [27]byte

func (h *Header) String() string {
	version := binary.LittleEndian.Uint32(h[23:27])
	return fmt.Sprint(string(h[0:20]), version)
}
