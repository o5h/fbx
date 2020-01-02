package fbx

import (
	"encoding/binary"
	"fmt"
)

type Header [27]byte

func (h Header) String() string {
	return fmt.Sprint(string(h[0:20]), h.Version())
}

func (h Header) Version() uint32 {
	return binary.LittleEndian.Uint32(h[23:27])
}
