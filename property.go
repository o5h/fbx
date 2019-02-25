package fbx

import (
	"fmt"
)

type Property struct {
	TypeCode byte
	Data     interface{}
}

func (p *Property) AsString() string {
	return p.Data.(string)
}

func (p *Property) AsBytes() []byte {
	return p.Data.([]byte)
}

func (p *Property) AsInt8() int8 {
	return p.Data.(int8)
}

func (p *Property) AsInt16() int16 {
	return p.Data.(int16)
}

func (p *Property) AsInt32() int32 {
	return p.Data.(int32)
}

func (p *Property) AsInt64() int64 {
	return p.Data.(int64)
}

func (p *Property) AsFloat32() float32 {
	return p.Data.(float32)
}

func (p *Property) AsFloat64() float64 {
	return p.Data.(float64)
}

func (p *Property) AsFloat32Slice() []float32 {
	return p.Data.(*Array).Data.([]float32)
}

func (p *Property) AsFloat64Slice() []float64 {
	return p.Data.(*Array).Data.([]float64)
}

func (p *Property) AsInt32Slice() []int32 {
	return p.Data.(*Array).Data.([]int32)
}

func (p *Property) AsInt64Slice() []int64 {
	return p.Data.(*Array).Data.([]int64)
}

func (p *Property) String() string {
	return fmt.Sprintf("%v", p.Data)
}
