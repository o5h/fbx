package fbx

import (
	"fmt"
)

type PropertyType byte

func (t PropertyType) String() string {
	return string([]byte{byte(t)})
}

type Property struct {
	TypeCode PropertyType
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

func (p *Property) AsFloat32Slice() (a []float32, ok bool) {
	a, ok = p.Data.(*Array).Data.([]float32)
	return
}

func (p *Property) AsFloat64Slice() (a []float64, ok bool) {
	a, ok = p.Data.(*Array).Data.([]float64)
	return
}

func (p *Property) AsInt32Slice() (a []int32, ok bool) {
	a, ok = p.Data.(*Array).Data.([]int32)
	return
}

func (p *Property) AsInt64Slice() (a []int64, ok bool) {
	a, ok = p.Data.(*Array).Data.([]int64)
	return
}

func (p *Property) String() string {
	return fmt.Sprintf("'%v:%v'", p.TypeCode, p.Data)
}
