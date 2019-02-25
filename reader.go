package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type FBXReader struct {
	FBX      *FBX
	Position int64
	Error    error
}

func NewReader() *FBXReader {
	return &FBXReader{&FBX{}, 0, nil}
}

func (fr *FBXReader) ReadFrom(r io.Reader) (n int64, err error) {
	fr.FBX.Header = fr.ReadHeaderFrom(r)
	if err != nil {
		return
	}

	fr.FBX.Top = fr.ReadNodeFrom(r)
	if fr.Error != nil {
		return
	}

	for {
		node := fr.ReadNodeFrom(r)
		if fr.Error != nil {
			break
		}
		if node.IsEmpty() {
			break
		}
		fr.FBX.Nodes = append(fr.FBX.Nodes, node)
	}

	return
}

func (fr *FBXReader) ReadHeaderFrom(r io.Reader) (header *Header) {
	header = &Header{}
	var i int
	i, fr.Error = r.Read(header[:])
	fr.Position += int64(i)
	return
}

func (fr *FBXReader) ReadNodeFrom(r io.Reader) (node *Node) {
	node = &Node{}

	node.EndOffset = fr.readUint32(r)
	if fr.Error != nil {
		return
	}

	node.NumProperties = fr.readUint32(r)
	if fr.Error != nil {
		return
	}

	node.PropertyListLen = fr.readUint32(r)
	if fr.Error != nil {
		return
	}

	node.NameLen = fr.readUint8(r)
	if fr.Error != nil {
		return
	}

	bb := make([]byte, node.NameLen)
	var i int
	i, fr.Error = io.ReadFull(r, bb)
	if fr.Error != nil {
		return
	}
	node.Name = string(bb)
	fr.Position += int64(i)

	if node.IsEmpty() {
		return
	}

	for np := uint32(0); np < node.NumProperties; np++ {
		p := fr.ReadPropertyFrom(r)
		if fr.Error != nil {
			return
		}
		node.Properties = append(node.Properties, p)
	}

	for {
		if fr.Position >= int64(node.EndOffset) {
			break
		}

		subNode := fr.ReadNodeFrom(r)
		if fr.Error != nil {
			break
		}
		if subNode.IsEmpty() {
			break
		}
		node.NestedNodes = append(node.NestedNodes, subNode)
	}

	return node
}

func (fr *FBXReader) ReadPropertyFrom(r io.Reader) (p *Property) {
	var nn int64
	p = &Property{}
	p.TypeCode = fr.readUint8(r)

	switch p.TypeCode {
	case 'S':
		p.Data = fr.readString(r)
	case 'R':
		p.Data = fr.readBytes(r)
	case 'Y':
		p.Data = fr.readInt16(r)
	case 'C':
		p.Data = fr.readInt8(r) != 0
	case 'I':
		p.Data = fr.readInt32(r)
	case 'F':
		p.Data = fr.readFloat32(r)
	case 'D':
		p.Data = fr.readFloat64(r)
	case 'L':
		p.Data = fr.readInt64(r)
	case 'f':
		p.Data = fr.readArray(r, 4,
			func(len uint32) interface{} {
				data := make([]float32, len)
				return data
			})
	case 'd':
		p.Data = fr.readArray(r, 8,
			func(len uint32) interface{} {
				data := make([]float64, len)
				return data
			})
	case 'i':
		p.Data = fr.readArray(r, 4,
			func(len uint32) interface{} {
				data := make([]int32, len)
				return data
			})
	case 'l':
		p.Data = fr.readArray(r, 8,
			func(len uint32) interface{} {
				data := make([]int64, len)
				return data
			})
	case 'b':
		var tmp []byte
		array := fr.readArray(r, 1,
			func(len uint32) interface{} {
				tmp = make([]byte, len)
				return tmp
			})
		data := make([]bool, len(tmp))
		for i, b := range tmp {
			data[i] = (b == 1)
		}
		array.Data = tmp
		p.Data = array
	default:
		panic(fmt.Sprintf("unsupported type '%s'", string(p.TypeCode)))
	}
	fr.Position += nn
	return
}

func (fr *FBXReader) readArrayHeader(r io.Reader, a *Array) {
	a.ArrayLength = fr.readUint32(r)
	if fr.Error != nil {
		return
	}
	a.Encoding = fr.readUint32(r)
	if fr.Error != nil {
		return
	}

	a.CompressedLength = fr.readUint32(r)
	if fr.Error != nil {
		return
	}
	return
}

func (fr *FBXReader) readArray(r io.Reader, size uint32, slicer func(len uint32) interface{}) *Array {
	a := &Array{}
	fr.readArrayHeader(r, a)
	if fr.Error != nil {
		return nil
	}
	data := slicer(a.ArrayLength)
	var nn int64
	nn, fr.Error = readArrayData(r, size, a, data)
	a.Data = data

	fr.Position += nn
	return a
}

func readArrayData(r io.Reader, size uint32, a *Array, data interface{}) (n int64, err error) {
	if a.Encoding == 0 {
		err = binary.Read(r, binary.LittleEndian, data)
		if err != nil {
			return
		}
		n += int64(size * a.ArrayLength)
	} else {
		var compressedBytes = make([]byte, a.CompressedLength)
		err = binary.Read(r, binary.LittleEndian, &compressedBytes)
		if err != nil {
			return
		}
		n += int64(a.CompressedLength)
		err = uncompress(compressedBytes, data)
	}
	return
}

func uncompress(b []byte, data interface{}) error {
	buf := bytes.NewBuffer(b)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return err
	}
	defer r.Close()
	err = binary.Read(r, binary.LittleEndian, data)
	return err
}

func (fr *FBXReader) readUint32(r io.Reader) uint32 {
	var data uint32
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 4
	return data
}

func (fr *FBXReader) readUint8(r io.Reader) uint8 {
	var data uint8
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position++
	return data
}

func (fr *FBXReader) readInt16(r io.Reader) int16 {
	var data int16
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 2
	return data
}

func (fr *FBXReader) readInt32(r io.Reader) int32 {
	var data int32
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 4
	return data
}

func (fr *FBXReader) readFloat32(r io.Reader) float32 {
	var data float32
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 4
	return data
}

func (fr *FBXReader) readFloat64(r io.Reader) float64 {
	var data float64
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 8
	return data
}

func (fr *FBXReader) readInt64(r io.Reader) int64 {
	var data int64
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position += 8
	return data
}

func (fr *FBXReader) readInt8(r io.Reader) int8 {
	var data int8
	fr.Error = binary.Read(r, binary.LittleEndian, &data)
	fr.Position++
	return data
}

func (fr *FBXReader) readString(r io.Reader) string {
	len := fr.readUint32(r)
	if fr.Error != nil {
		return ""
	}
	bb := make([]byte, len)
	var i int
	i, fr.Error = io.ReadFull(r, bb)
	if fr.Error != nil {
		return ""
	}
	fr.Position += int64(i)
	return string(bb)
}

func (fr *FBXReader) readBytes(r io.Reader) []byte {
	len := fr.readUint32(r)
	if fr.Error != nil {
		return nil
	}

	bb := make([]byte, len)
	var i int
	i, fr.Error = io.ReadFull(r, bb)
	if fr.Error != nil {
		return nil
	}
	fr.Position += int64(i)
	return bb
}

func ReadFrom(r io.Reader) (*FBX, error) {
	reader := NewReader()
	reader.ReadFrom(r)
	return reader.FBX, reader.Error
}
