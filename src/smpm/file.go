// SMPM - Provide Minecraft Multiplayer Map Chunk file format parsing capabilities
package smpm

import (
	"encoding/binary"
	"fmt"
	"io"
)

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
		return nil
	}

	total := f.metadata.ChunkColumnCount()
	for i := int16(0); i < total; i++ {
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
		err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
			return c.Metadata
		})
		if err != nil {
			return
		}
		err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
			return c.Light
		})
		if err != nil {
			return
		}
		if f.metadata.HasSkylightData() {
			err = eachChunkRead(chunks, metadata.PrimaryBitmap, func(c *Chunk) []byte {
				return c.Skylight
			})
			if err != nil {
				return
			}
		}
		err = eachChunkRead(chunks, metadata.AddBitmap, func(c *Chunk) []byte {
			return c.Add
		})
		if err != nil {
			return
		}
		if f.metadata.IsGroundUpContinuous() {
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

func (f *File) readBytes(bytes []byte) error {
	return binary.Read(f.reader, f.ByteOrder, &bytes)
}
