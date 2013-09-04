// SMPM - Provide Minecraft Multiplayer Map Chunk file format parsing capabilities
package smpm

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ChunkSize       = 16 * 16 * 16 // blocks
	ChunksPerColumn = 16           // number of chunks per column
	ChunkBiomeSize  = 256
	halfDataDivisor = 2
	HighBits        = 0xf0 // high ordered bit-mask of half-byte stores
	LowBits         = 0x0f // low ordered bit-mask of half-byte stores
)

type Metadata interface {
	ChunkColumnCount() int16
	NextMetadata() ChunkColumnMetadata
	HasSkylightData() bool
	IsGroundUpContinuous() bool
}

type ColumnPoint struct {
	X, Z int32
}

type ChunkColumnMetadata struct {
	X, Z                     int32
	PrimaryBitmap, AddBitmap uint16
}

// represents a 16x16x16 blocks
// order Y, Z, X ascending (0 -> 15)
type Chunk struct {
	Types    []byte
	Metadata []byte
	Light    []byte
	Skylight []byte
	Add      []byte
}

func NewChunk() *Chunk {
	return &Chunk{
		Types:    make([]byte, ChunkSize),
		Metadata: make([]byte, ChunkSize/halfDataDivisor),
		Light:    make([]byte, ChunkSize/halfDataDivisor),
		Skylight: make([]byte, ChunkSize/halfDataDivisor),
		Add:      make([]byte, ChunkSize/halfDataDivisor),
	}
}

func NewChunkSlice(size int) []*Chunk {
	chunks := make([]*Chunk, size)
	for i := 0; i < size; i++ {
		chunks[i] = NewChunk()
	}
	return chunks
}

// represents 16x256x16 blocks
type ChunkColumn struct {
	Chunks   []*Chunk             // order of chunks are from bottom to top
	Biome    [ChunkBiomeSize]byte // 16x16 of the biome for each X, Z coordinate
	Metadata *ChunkColumnMetadata
}

type File struct {
	reader    io.Reader
	metadata  Metadata
	ByteOrder binary.ByteOrder
}

func NewFile(reader io.Reader, metadata Metadata) *File {
	return &File{
		reader:    reader,
		metadata:  metadata,
		ByteOrder: binary.BigEndian,
	}
}

func (f *File) Parse() (columns []ChunkColumn, err error) {
	columns = make([]ChunkColumn, 0, f.metadata.ChunkColumnCount())

	if f.reader == nil {
		err = fmt.Errorf("No Reader specified")
		return
	}

	eachChunkRead := func(chunks []*Chunk, mask uint16, read func(chunk *Chunk) []byte) error {
		blocksToRead := 0
		for i, chunk := range chunks {
			if (mask & (1 << uint16(i))) > 0 {
				bytes := read(chunk)
				err := f.readBytes(bytes)
				blocksToRead += 1
				if err != nil {
					return err
				}
			}
		}
		fmt.Printf("    Blocks read: %d / %d\n", blocksToRead, len(chunks))
		return nil
	}

	total := f.metadata.ChunkColumnCount()
	fmt.Printf("Total Chunks: #%d\n", total)
	for i := int16(0); i < total; i++ {
		fmt.Printf("Chunk: #%d\n", i+1)
		fmt.Printf(" Types\n")
		metadata := f.metadata.NextMetadata()
		column := ChunkColumn{
			Chunks:   NewChunkSlice(ChunksPerColumn),
			Metadata: &metadata,
		}
		chunks := column.Chunks
		err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
			return c.Types
		})
		if err != nil {
			return
		}
		fmt.Printf(" Metadata\n")
		err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
			return c.Metadata
		})
		if err != nil {
			return
		}
		fmt.Printf(" Light\n")
		err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
			return c.Light
		})
		if err != nil {
			return
		}
		if f.metadata.HasSkylightData() {
			fmt.Printf(" SkyLight\n")
			err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
				return c.Skylight
			})
			if err != nil {
				return
			}
		}
		fmt.Printf(" Add\n")
		err = eachChunkRead(chunks, metadata.AddBitmap, func(c *Chunk) []byte {
			return c.Add
		})
		if err != nil {
			return
		}
		if f.metadata.IsGroundUpContinuous() {
			fmt.Printf(" Biome\n")
			err = f.readBytes(column.Biome[:])
			if err != nil {
				return
			}
		}
		//columns[i] = column
		columns = append(columns, column)
	}
	return
}

var count = 0

func hash(bytes []byte, prefix string) {
	hasher := md5.New()
	hasher.Write(bytes)
	fmt.Printf(" ==> %s %x\n", prefix, hasher.Sum(nil))
}

func (f *File) readBytes(bytes []byte) error {
	defer func() {
		hasher := md5.New()
		hasher.Write(bytes)
		count += len(bytes)
		fmt.Printf("    Read %d bytes (t=%d, hash=%x)\n", len(bytes), count, hasher.Sum(nil))
	}()
	return binary.Read(f.reader, f.ByteOrder, &bytes)
}
