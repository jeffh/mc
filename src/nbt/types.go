package nbt

type TagType byte

const (
    TagTypeEnd TagType = iota
    TagTypeByte
    TagTypeShort
    TagTypeInt
    TagTypeLong
    TagTypeFloat
    TagTypeDouble
    TagTypeByteArray
    TagTypeString
    TagTypeList
    TagTypeCompound
    TagTypeIntArray
    TagTypeInvalid = 0xff
)

type Tag struct {
    Name string
    Type TagType
    Value interface{}
}

var InvalidTag = Tag{Type: TagTypeInvalid}

//type Byte byte
//type Short int16
//type Int int32
//type Long int64
////type Float float32
//type Double float64
//type ByteArray []byte
//type String string
type List []Tag
type Compound map[string]Tag
type IntArray []int32
