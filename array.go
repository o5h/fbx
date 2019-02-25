package fbx

import (
	"fmt"
)

type Array struct {
	ArrayLength      uint32
	Encoding         uint32
	CompressedLength uint32
	Data             interface{}
}

func (a *Array) String() string {
	return fmt.Sprintf("%v", a.Data)
}
