package fbx_test

import (
	"io"
	"os"
	"testing"

	"github.com/o5h/fbx"
	"github.com/o5h/testing/assert"
)

func TestRead(t *testing.T) {
	f, _ := os.Open("testdata/cube.fbx")

	defer f.Close()

	reader := fbx.NewReader()
	reader.ReadFrom(f)

	cur, _ := f.Seek(0, io.SeekCurrent)
	assert.Eq(t, reader.Position, cur)

	fbxData := reader.FBX
	ibo := fbxData.Filter(fbx.FilterName("PolygonVertexIndex"))[0]
	iboData := ibo.Properties[0].AsInt32Slice()

	t.Log(fbxData)

	assert.Eq(t, iboData, []int32{
		0, 2, -4,
		7, 5, -5,
		4, 1, -1,
		5, 2, -2,
		2, 7, -4,
		0, 7, -5,
		0, 1, -3,
		7, 6, -6,
		4, 5, -2,
		5, 6, -3,
		2, 6, -8,
		0, 3, -8})

	vbo := fbxData.Filter(fbx.FilterName("Vertices"))[0]
	vboData := vbo.Properties[0].AsFloat64Slice()

	assert.EqSlice(t, vboData, []float64{
		1, 0.999999940395355, -1,
		1, -1, -1,
		-1.00000011920929, -0.999999821186066,
		-1, -0.999999642372131, 1.00000035762787,
		-1, 1.00000047683716, 0.999999463558197,
		1, 0.999999344348907, -1.00000059604645,
		1, -1.00000035762787, -0.999999642372131,
		1, -0.999999940395355, 1, 1},
		0.0000001)
}
