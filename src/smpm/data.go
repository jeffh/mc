package smpm

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
