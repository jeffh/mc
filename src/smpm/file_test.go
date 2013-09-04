package smpm

import (
	"bytes"
	"compress/zlib"
	. "github.com/jeffh/goexpect"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type sampleMetadata struct {
	hasSkylightData      bool
	isGroundUpContinuous bool
	metadatas            []ChunkColumnMetadata
	index                int16
}

func (s *sampleMetadata) ChunkColumnCount() int16 {
	return int16(len(s.metadatas))
}

func (s *sampleMetadata) HasSkylightData() bool {
	return s.hasSkylightData
}

func (s *sampleMetadata) IsGroundUpContinuous() bool {
	return s.isGroundUpContinuous
}

func (s *sampleMetadata) NextMetadata() ChunkColumnMetadata {
	index := s.index
	s.index += 1
	return s.metadatas[index]
}

var sampleChunkColumnMetadata = &sampleMetadata{
	hasSkylightData:      true,
	isGroundUpContinuous: true,
	metadatas: []ChunkColumnMetadata{
		{-17, 14, 31, 0},
		{-16, 14, 31, 0},
		{-16, 15, 31, 0},
		{-17, 15, 31, 0},
		{-18, 15, 31, 0},
	},
}

func fixtureAsReader(filename string) io.Reader {
	buf := bytes.NewBuffer([]byte{})
	filepath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	fi, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	_, err = io.Copy(buf, fi)
	if err != nil {
		panic(err)
	}
	return buf
}

func TestReadingWithNoReaderIsAnError(t *testing.T) {
	smpm := NewFile(nil, sampleChunkColumnMetadata)
	_, err := smpm.Parse()
	Expect(t, err.Error(), ToEqual, "No Reader specified")
}

func fill(array []byte, value byte) {
	for i := 0; i < len(array); i++ {
		array[i] = value
	}
}

func TestReadingBulkChunk(t *testing.T) {
	r := fixtureAsReader("fixtures/sample.bin.zlib")
	reader, err := zlib.NewReader(r)
	Expect(t, err, ToBeNil)

	smpm := NewFile(reader, sampleChunkColumnMetadata)
	columns, err := smpm.Parse()
	Expect(t, err, ToBeNil)

	Expect(t, columns, ToBeLengthOf, 5)
	Expect(t, columns[0].Metadata, ToEqual, &sampleChunkColumnMetadata.metadatas[0])
	Expect(t, columns[1].Metadata, ToEqual, &sampleChunkColumnMetadata.metadatas[1])
	Expect(t, columns[2].Metadata, ToEqual, &sampleChunkColumnMetadata.metadatas[2])
	Expect(t, columns[3].Metadata, ToEqual, &sampleChunkColumnMetadata.metadatas[3])
	Expect(t, columns[4].Metadata, ToEqual, &sampleChunkColumnMetadata.metadatas[4])

	chunk := columns[0].Chunks[0]

	Expect(t, chunk.Types, ToEqual, types)
	Expect(t, chunk.Light, ToEqual, light)
	Expect(t, chunk.Metadata, ToEqual, make([]byte, 2048))
	Expect(t, chunk.Add, ToEqual, make([]byte, 2048))
	Expect(t, chunk.Skylight, ToEqual, make([]byte, 2048))
	biome := [256]byte{}
	fill(biome[:], 5)
	Expect(t, columns[0].Biome, ToEqual, biome)
}

var light = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0,
	0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 255, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 255, 255, 0, 0, 0, 0, 0, 255, 255, 255,
	255, 255, 15, 0, 0, 255, 255, 255, 255, 255, 15, 0, 0, 255, 255, 255, 255, 255,
	15, 0, 0, 255, 0, 240, 255, 255, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240, 255, 255,
	0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 15, 0, 255, 255, 255, 255, 255,
	255, 255, 15, 255, 255, 255, 255, 255, 255, 255, 255, 15, 0, 240, 255, 255, 255,
	255, 255, 0, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0,
	240, 255, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 240,
	255, 255, 255, 255, 255, 255, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 240, 255, 255, 255, 255, 0, 0, 0, 0,
	0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0,
	0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 240, 255, 255, 0, 0,
	0, 0, 0, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	240, 15, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0,
	0, 0, 0, 0, 240, 255, 255, 0, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 0, 255,
	0, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 255,
	255, 0, 0, 0, 0, 0, 240, 255, 255, 0, 0, 0, 0, 0, 255, 255, 255, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0,
	0, 255, 0, 0, 0, 0, 0, 0, 0, 255, 15, 0, 0, 0, 0, 0, 0, 255, 15, 0, 0, 0, 0,
	240, 255, 255, 15, 0, 0, 0, 0, 255, 255, 255, 0, 0, 0, 0, 0, 240, 255, 0, 0, 0,
	0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0,
	0, 0, 0, 240, 255, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 240, 255, 255, 0,
	240, 255, 0, 0, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 222, 0, 0, 0, 0, 0, 0, 0, 238, 13, 0, 0, 0, 0, 0, 0,
	238, 14, 144, 0, 0, 0, 0, 0, 238, 222, 0, 137, 0, 0, 0, 0, 238, 222, 188, 154,
	120, 0, 0, 0, 238, 13, 171, 137, 103, 5, 0, 224, 221, 0, 0, 120, 86, 6, 0, 224,
	0, 0, 0, 0, 101, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 238, 0, 0, 0, 0, 0, 0, 0, 238,
	160, 203, 0, 0, 0, 0, 224, 238, 176, 220, 0, 0, 0, 0, 238, 238, 192, 237, 238,
	0, 0, 224, 238, 238, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 0, 0,
	0, 0, 0, 0, 0, 205, 0, 0, 0, 0, 0, 0, 0, 221, 12, 0, 0, 0, 0, 0, 0, 221, 205,
	171, 0, 0, 0, 0, 0, 221, 205, 171, 137, 7, 0, 0, 0, 221, 205, 171, 137, 103, 5,
	0, 0, 221, 188, 154, 120, 86, 6, 0, 0, 204, 0, 128, 103, 101, 135, 0, 0, 187, 0,
	0, 80, 118, 152, 10, 220, 170, 9, 0, 0, 135, 169, 203, 221, 153, 169, 0, 0, 0,
	169, 203, 221, 152, 186, 0, 0, 0, 176, 220, 221, 169, 203, 0, 0, 0, 0, 0, 0,
	176, 220, 221, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	11, 0, 0, 0, 0, 0, 0, 0, 188, 10, 0, 0, 0, 0, 0, 0, 204, 171, 0, 0, 0, 0, 0, 0,
	204, 188, 154, 8, 0, 0, 0, 0, 204, 188, 154, 120, 6, 0, 0, 0, 204, 188, 154,
	120, 86, 0, 0, 0, 204, 171, 137, 103, 69, 101, 0, 0, 187, 0, 112, 86, 84, 118,
	8, 176, 170, 9, 0, 64, 101, 135, 169, 203, 153, 8, 0, 0, 118, 152, 186, 204,
	136, 152, 0, 0, 112, 152, 186, 204, 135, 169, 0, 0, 0, 169, 203, 204, 152, 186,
	0, 0, 0, 0, 0, 0, 160, 203, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 171, 9, 0, 0, 0, 0, 0, 0, 187, 154, 8, 0,
	0, 0, 0, 0, 187, 171, 137, 7, 0, 0, 0, 0, 187, 171, 137, 103, 0, 0, 0, 0, 187,
	171, 137, 103, 5, 0, 0, 0, 187, 10, 120, 86, 50, 4, 0, 0, 170, 0, 0, 37, 67,
	101, 0, 160, 153, 0, 0, 0, 84, 118, 152, 186, 136, 7, 0, 0, 96, 135, 169, 187,
	119, 7, 0, 0, 96, 135, 169, 187, 118, 152, 0, 0, 0, 144, 186, 187, 128, 169, 0,
	0, 0, 0, 0, 0, 144, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 170, 128, 7, 0, 0, 0,
	0, 0, 170, 154, 8, 0, 0, 0, 0, 0, 170, 10, 120, 0, 0, 0, 0, 0, 170, 10, 0, 0, 0,
	0, 0, 0, 170, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 136, 0, 0, 0, 0, 101,
	135, 160, 119, 0, 0, 0, 0, 118, 152, 10, 102, 0, 0, 0, 0, 112, 152, 170, 0, 0,
	0, 0, 0, 0, 160, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var types = []byte{
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 15, 1, 7, 7, 7, 7, 7, 7,
	1, 7, 7, 1, 7, 7, 1, 7, 1, 1, 7, 7, 7, 7, 1, 1, 1, 7, 1, 7, 7, 7, 7, 7, 1, 7, 7,
	7, 7, 1, 1, 1, 1, 7, 7, 1, 1, 7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 7, 1, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 1, 7, 7, 1, 1, 7, 1, 7, 7, 7, 1, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 7, 1, 7, 1, 1, 7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 7, 7,
	7, 7, 7, 1, 1, 7, 7, 7, 1, 7, 7, 7, 7, 1, 1, 7, 7, 7, 7, 1, 7, 7, 7, 7, 7, 1, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 7, 1, 1, 7, 7,
	7, 7, 1, 7, 7, 7, 7, 7, 7, 7, 1, 1, 7, 7, 7, 3, 7, 1, 7, 7, 7, 7, 1, 7, 7, 7, 7,
	7, 7, 7, 7, 3, 7, 7, 7, 1, 1, 7, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 1, 7, 7, 1, 7, 7, 1, 7, 1, 15, 7, 7, 7, 3, 3, 7, 7, 1, 1, 7,
	1, 1, 1, 7, 7, 1, 7, 7, 1, 1, 7, 1, 7, 7, 7, 7, 7, 7, 7, 7, 1, 1, 7, 7, 1, 1, 7,
	7, 7, 1, 7, 1, 7, 1, 1, 1, 7, 1, 1, 7, 7, 7, 1, 1, 7, 7, 1, 7, 7, 1, 7, 1, 1, 1,
	7, 7, 7, 7, 1, 1, 7, 7, 1, 7, 7, 1, 7, 1, 7, 1, 7, 7, 7, 1, 7, 7, 7, 1, 1, 1, 1,
	1, 7, 7, 7, 1, 7, 7, 1, 7, 7, 7, 1, 1, 7, 7, 1, 7, 1, 1, 7, 7, 7, 7, 7, 7, 1, 1,
	7, 7, 1, 7, 7, 1, 1, 1, 7, 7, 7, 7, 1, 7, 7, 7, 7, 1, 7, 1, 1, 1, 1, 1, 1, 7, 7,
	1, 7, 7, 1, 1, 7, 7, 7, 7, 7, 1, 7, 1, 7, 1, 1, 1, 1, 7, 1, 7, 7, 1, 7, 7, 7, 1,
	73, 7, 1, 7, 7, 7, 7, 7, 7, 7, 7, 1, 3, 7, 7, 7, 1, 1, 7, 1, 7, 1, 7, 1, 1, 7,
	7, 1, 7, 7, 7, 1, 7, 7, 7, 7, 7, 1, 1, 7, 7, 1, 7, 7, 7, 7, 7, 7, 7, 1, 1, 1, 1,
	7, 7, 7, 1, 7, 7, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 7, 7, 7, 7, 7, 1, 7, 7, 7, 1,
	7, 1, 1, 1, 7, 7, 7, 3, 7, 3, 3, 1, 1, 1, 7, 7, 7, 1, 1, 1, 7, 1, 1, 1, 1, 1, 7,
	7, 1, 7, 7, 1, 1, 1, 7, 1, 7, 7, 7, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 7, 1, 7, 7, 7,
	1, 7, 7, 7, 1, 1, 7, 1, 1, 1, 7, 7, 7, 1, 1, 7, 1, 1, 7, 7, 1, 7, 1, 1, 1, 7, 7,
	7, 1, 7, 7, 7, 1, 7, 7, 1, 7, 7, 7, 1, 7, 1, 7, 1, 7, 7, 1, 7, 7, 1, 7, 7, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 7, 7, 7, 1, 7, 7, 1, 1, 1, 7, 7, 1, 7, 7, 1, 7, 7, 7, 1, 1,
	7, 1, 7, 7, 1, 1, 1, 1, 73, 73, 7, 1, 1, 7, 1, 1, 1, 1, 7, 7, 7, 7, 7, 1, 1, 1,
	7, 1, 1, 1, 1, 1, 7, 1, 1, 7, 1, 7, 1, 1, 1, 1, 1, 7, 1, 1, 7, 1, 1, 1, 7, 1, 1,
	7, 1, 1, 1, 1, 1, 1, 1, 1, 7, 1, 1, 1, 1, 1, 7, 3, 3, 3, 1, 7, 7, 7, 7, 1, 1, 1,
	1, 1, 1, 1, 1, 7, 7, 3, 15, 7, 1, 1, 1, 7, 1, 1, 7, 1, 7, 1, 1, 1, 7, 7, 7, 1,
	1, 7, 7, 1, 1, 7, 1, 7, 1, 7, 1, 1, 1, 1, 7, 1, 1, 1, 1, 7, 1, 1, 1, 3, 1, 1, 1,
	1, 7, 1, 1, 1, 1, 1, 15, 1, 1, 1, 1, 1, 1, 1, 1, 7, 7, 1, 1, 1, 73, 1, 1, 7, 1,
	1, 1, 1, 1, 7, 1, 1, 7, 1, 1, 1, 73, 1, 1, 7, 1, 1, 7, 1, 1, 7, 1, 7, 1, 1, 1,
	1, 1, 1, 1, 7, 1, 1, 1, 7, 1, 1, 1, 1, 1, 1, 1, 1, 7, 1, 1, 1, 7, 7, 7, 1, 1, 1,
	7, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 7, 7, 7, 7, 1, 1, 1, 1, 7, 1, 1, 7,
	1, 1, 1, 1, 1, 1, 7, 1, 1, 1, 73, 7, 1, 1, 1, 1, 1, 1, 7, 7, 1, 1, 1, 7, 1, 1,
	1, 1, 1, 1, 1, 7, 1, 1, 7, 1, 1, 1, 1, 1, 1, 7, 1, 1, 7, 1, 7, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 7, 1, 1, 1, 7, 1, 1, 1, 1, 1, 1, 1, 7, 7, 1, 7, 1, 1, 1, 1, 7, 1, 1,
	1, 7, 1, 1, 1, 7, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 7, 1, 1, 7, 7, 1, 1, 3, 15, 15,
	1, 1, 1, 1, 1, 1, 7, 1, 1, 1, 1, 1, 1, 1, 15, 7, 1, 7, 1, 1, 1, 7, 1, 1, 1, 7,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 15, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 15,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 73, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 11, 10, 11, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 73, 11, 11, 11, 11, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 14, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11,
	10, 11, 11, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 11, 11, 11, 11, 11, 11, 11, 11,
	11, 10, 11, 1, 1, 1, 1, 1, 11, 11, 10, 11, 11, 11, 11, 11, 10, 11, 11, 1, 1, 1,
	1, 1, 11, 11, 1, 1, 1, 11, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 73, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 14, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 10, 10, 11, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 11, 11, 11, 11, 11, 11, 11, 11, 10, 11, 11, 11, 11, 1, 1, 1, 11, 10, 11, 10,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 1, 11, 11, 11, 11, 11, 11, 11, 11,
	10, 11, 11, 11, 11, 11, 11, 11, 11, 1, 1, 1, 73, 11, 11, 10, 11, 11, 11, 10, 11,
	11, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11, 10, 1, 1, 1, 1, 1, 1,
	1, 1, 73, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 73, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 11, 11, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11, 11, 1, 1, 1,
	1, 1, 1, 1, 14, 73, 1, 1, 1, 1, 1, 10, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11, 11, 11, 11, 11, 11, 11, 11,
	11, 10, 11, 11, 1, 1, 11, 11, 11, 11, 11, 11, 11, 10, 11, 11, 11, 11, 11, 11,
	10, 11, 10, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 10, 1, 1, 1,
	1, 1, 1, 1, 10, 11, 11, 11, 11, 10, 11, 11, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 11, 11, 11, 10, 1, 1, 1, 1, 1, 1, 1, 73, 1, 1, 1, 1, 1, 1, 11, 11, 1, 1, 1,
	1, 1, 1, 1, 73, 1, 1, 1, 1, 1, 1, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	11, 10, 11, 15, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 11, 1, 1, 1, 1, 1,
	1, 1, 1, 73, 73, 1, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 73, 73, 11, 11,
	11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 11, 11, 11, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 11,
	11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 11, 11, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 10, 11, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 11, 11, 11, 15, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11,
	11, 15, 1, 1, 1, 1, 1, 1, 1, 1, 73, 1, 11, 11, 10, 11, 11, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 11, 11, 11, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 10, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 11, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 10, 11, 10, 11, 10, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 11, 11, 10, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10, 11,
	11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 16, 16, 1, 1, 1, 11, 11, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 16, 16, 1, 1, 11, 11, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 11, 11, 10, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 11, 11, 11, 11, 11, 1, 1,
	1, 11, 11, 11, 1, 1, 1, 1, 11, 10, 11, 10, 11, 11, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 13,
	13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 13, 13, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 13, 13, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,
	1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 16, 16, 1, 1, 1, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 16, 16, 16, 1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 16,
	16, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 13, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 13, 13, 13, 13,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 13, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 0, 0, 0, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 13,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
	1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0,
	0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 1, 16, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 16, 16, 16, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 16, 16, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 13, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 13, 13, 13, 13, 13, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0,
	1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 16, 16, 15, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 1, 1, 16, 16, 15, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 15, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 13, 13,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 13, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 13, 13, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 13,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
	1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0,
	0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 15, 15, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1,
	1, 1, 15, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 13, 13, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 13, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
}