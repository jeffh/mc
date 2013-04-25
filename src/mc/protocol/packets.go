package protocol

import (
	"bytes"
	"compress/gzip"
	"io"
)

const Version = 60 // minecraft protocol version supported

///////////////////////////////////////////////////////

var basePacketMapper *StdPacketMapper
var ClientPacketMapper, ServerPacketMapper *StdPacketMapper

// sets up the mapping of opcodes to structs
func init() {
	basePacketMapper = NewStdPacketMapper(nil)
	basePacketMapper.Define(0x00, KeepAlive{})
	basePacketMapper.Define(0x01, LoginRequest{})
	basePacketMapper.Define(0x02, Handshake{})
	basePacketMapper.Define(0x03, ChatMessage{})
	basePacketMapper.Define(0x04, TimeUpdate{})
	basePacketMapper.Define(0x05, EntityEquipment{})
	basePacketMapper.Define(0x06, SpawnPosition{})
	basePacketMapper.Define(0x07, UseEntity{})
	basePacketMapper.Define(0x08, UpdateHealth{})
	basePacketMapper.Define(0x09, Respawn{})
	basePacketMapper.Define(0x0A, Player{})
	basePacketMapper.Define(0x0B, PlayerPosition{})
	basePacketMapper.Define(0x0C, PlayerLook{})
	basePacketMapper.Define(0x0E, PlayerDigging{})
	basePacketMapper.Define(0x0F, PlayerBlockPlacement{})
	basePacketMapper.Define(0x10, HeldItemChange{})
	basePacketMapper.Define(0x11, UseBed{})
	basePacketMapper.Define(0x12, Animation{})
	basePacketMapper.Define(0x13, EntityAction{})
	basePacketMapper.Define(0x14, SpawnNamedEntity{})
	basePacketMapper.Define(0x15, SpawnDroppedItem{})
	basePacketMapper.Define(0x16, CollectItem{})
	basePacketMapper.Define(0x17, SpawnObject{})
	basePacketMapper.Define(0x18, SpawnMob{})
	basePacketMapper.Define(0x19, SpawnPainting{})
	basePacketMapper.Define(0x1A, SpawnExperienceOrb{})
	basePacketMapper.Define(0x1C, EntityVelocity{})
	basePacketMapper.Define(0x1D, DestroyEntity{})
	basePacketMapper.Define(0x1E, CreateEntity{})
	basePacketMapper.Define(0x1F, EntityRelativeMove{})
	basePacketMapper.Define(0x20, EntityLook{})
	basePacketMapper.Define(0x21, EntityLookRelativeMove{})
	basePacketMapper.Define(0x22, EntityTeleport{})
	basePacketMapper.Define(0x23, EntityHeadLook{})
	basePacketMapper.Define(0x26, EntityStatus{})
	basePacketMapper.Define(0x27, AttachEntity{})
	basePacketMapper.Define(0x28, SetEntityMetadata{})
	basePacketMapper.Define(0x29, EntityEffect{})
	basePacketMapper.Define(0x2A, RemoveEntityEffect{})
	basePacketMapper.Define(0x2B, SetExperience{})
	basePacketMapper.Define(0x33, ChunkData{})
	basePacketMapper.Define(0x34, MultiBlockChange{})
	basePacketMapper.Define(0x35, BlockChange{})
	basePacketMapper.Define(0x36, BlockAction{})
	basePacketMapper.Define(0x37, BlockBreakAnimation{})
	basePacketMapper.Define(0x38, MapChunkBulk{})
	basePacketMapper.Define(0x3C, Explosion{})
	basePacketMapper.Define(0x3D, Effect{})
	basePacketMapper.Define(0x3E, NamedSoundEffect{})
	basePacketMapper.Define(0x46, ChangeGameState{})
	basePacketMapper.Define(0x47, GlobalEntity{})
	basePacketMapper.Define(0x64, OpenWindow{})
	basePacketMapper.Define(0x65, CloseWindow{})
	basePacketMapper.Define(0x66, ClickWindow{})
	basePacketMapper.Define(0x67, SetSlot{})
	basePacketMapper.Define(0x68, SetWindowItems{})
	basePacketMapper.Define(0x69, UpdateWindowProperty{})
	basePacketMapper.Define(0x6A, ConfirmTransaction{})
	basePacketMapper.Define(0x6B, CreativeInventoryAction{})
	basePacketMapper.Define(0x6C, EnchantItem{})
	basePacketMapper.Define(0x82, UpdateSign{})
	basePacketMapper.Define(0x83, ItemData{})
	basePacketMapper.Define(0x84, UpdateTileEntity{})
	basePacketMapper.Define(0xC8, IncrementStatistic{})
	basePacketMapper.Define(0xC9, PlayerListItem{})
	basePacketMapper.Define(0xCA, PlayerAbilities{})
	basePacketMapper.Define(0xCB, TabComplete{})
	basePacketMapper.Define(0xCC, ClientSettings{})
	basePacketMapper.Define(0xCD, ClientStatus{})
	basePacketMapper.Define(0xCE, ScoreboardObjective{})
	basePacketMapper.Define(0xCF, UpdateScore{})
	basePacketMapper.Define(0xD0, DisplayScoreboard{})
	basePacketMapper.Define(0xD1, Teams{})
	basePacketMapper.Define(0xFA, PluginMessage{})
	basePacketMapper.Define(0xFC, EncryptionKeyResponse{})
	basePacketMapper.Define(0xFD, EncryptionKeyRequest{})
	basePacketMapper.Define(0xFE, ServerListPing{})
	basePacketMapper.Define(0xFF, Disconnect{})
	// complicated b/c there's different formats for server/client
	//basePacketMapper.Define(0x0D, PlayerPositionLook{})

	// PlayerPositionLook has different ordering a fields for server and
	// clients... we handle the cases here.
	ServerPacketMapper = NewStdPacketMapper(basePacketMapper)
	ServerPacketMapper.DefineIncoming(0x0D, PlayerPositionLookForServer{})
	ServerPacketMapper.DefineOutgoing(0x0D, PlayerPositionLookForClient{})

	ClientPacketMapper = NewStdPacketMapper(basePacketMapper)
	ClientPacketMapper.DefineIncoming(0x0D, PlayerPositionLookForClient{})
	ClientPacketMapper.DefineOutgoing(0x0D, PlayerPositionLookForServer{})
}

///////////////////////////////////////////////////////

// All the data structures represented by the protocol are here.
// Remember that all java data types are signed.

///////////////////////////////////////////////////////

type KeepAlive struct {
	ID int32
}
type LoginRequest struct {
	EntityID   int32
	LevelType  string
	GameMode   GameMode
	Dimension  GameDimension
	Difficulty GameDifficulty
	NotUsed    int8
	MaxPlayers int8
}
type Handshake struct {
	Version  byte
	Username string
	Hostname string
	Port     int32
}
type ChatMessage struct {
	Message string // clients are limited to max of 100 characters
}
type TimeUpdate struct {
	WorldAge  int64
	TimeOfDay int64
}
type EntityEquipment struct {
	EntityID int32
	Slot     int16 // 0 = Held, 1-4 = Armor
	Item     Slot
}
type SpawnPosition struct {
	X, Y, Z int32
}
type UseEntity struct {
	User, Target      int32
	IsLeftMouseButton bool
}
type UpdateHealth struct {
	Health, Food int16
	Saturation   float32
}
type Respawn struct {
	Dimension   GameDimension
	Difficulty  GameDifficulty
	GameMode    GameMode
	WorldHeight int16
	LevelType   string
}
type Player struct {
	IsOnGround bool
}
type PlayerPosition struct {
	// ordering matters for parsing!
	X, Y, Stance, Z float64
	IsOnGround      bool
}
type PlayerLook struct {
	Yaw, Pitch float32
	IsOnGround bool
}

// it's slightly different based on who's sending it.
// see http://www.wiki.vg/Protocol#Player_Position_and_Look_.280x0D.29
type PlayerPositionLookForServer struct { // client -> server
	X, Y, Stance, Z float64
	Yaw, Pitch      float32
	IsOnGround      bool
}

// Converts this packet to the appropriate packet for sending to the Client.
func (p *PlayerPositionLookForServer) PacketForClient() *PlayerPositionLookForClient {
	return &PlayerPositionLookForClient{
		X:          p.X,
		Y:          p.Y,
		Stance:     p.Stance,
		Z:          p.Z,
		Yaw:        p.Yaw,
		Pitch:      p.Pitch,
		IsOnGround: p.IsOnGround,
	}
}

type PlayerPositionLookForClient struct { // server -> client
	X, Stance, Y, Z float64
	Yaw, Pitch      float32
	IsOnGround      bool
}

// Converts this packet to the appropriate packet for sending to the Server.
func (p *PlayerPositionLookForClient) PacketForServer() *PlayerPositionLookForServer {
	return &PlayerPositionLookForServer{
		X:          p.X,
		Y:          p.Y,
		Stance:     p.Stance,
		Z:          p.Z,
		Yaw:        p.Yaw,
		Pitch:      p.Pitch,
		IsOnGround: p.IsOnGround,
	}
}

type PlayerDigging struct {
	Status PlayerDiggingStatus
	X      int32
	Y      int8
	Z      int32
	Face   Face
}
type PlayerBlockPlacement struct {
	X         int32
	Y         byte
	Z         int32
	Direction int8
	ItemHeld  Slot
	CursorX   int8
	CursorY   int8
	CursorZ   int8
}
type HeldItemChange struct {
	SlotID int16 // 0-8
}
type UseBed struct { // sent by server to all nearby players
	EntityID int32
	Unknown  int8 // always zero?
	X        int32
	Y        int8
	Z        int32
}
type Animation struct {
	EntityID int32
	Type     AnimationType
}
type EntityAction struct {
	EntityID int32
	ActionID ActionType
}
type SpawnNamedEntity struct {
	EntityID    int32
	PlayerName  string
	X, Y, Z     int32
	Yaw, Pitch  int8
	CurrentItem int16
	Metadata    []EntityMetadata
}
type SpawnDroppedItem struct {
	EntityID              int32
	Slot                  Slot
	X, Y, Z               int32
	Rotation, Pitch, Roll int8
}
type CollectItem struct {
	CollectedEntityID int32
	ColloctorEntityID int32
}
type SpawnObject struct {
	EntityID                        int32
	Type                            EntityType
	X, Y, Z                         int32
	Data                            int32
	XVelocity, YVelocity, ZVelocity int16 // can be missing if data is 0
}
type SpawnMob struct {
	EntityID                        int32
	Type                            MobType
	X, Y, Z                         int32
	Pitch, HeadPitch, Yaw           int8
	XVelocity, YVelocity, ZVelocity int16
	Metadata                        []EntityMetadata
}
type SpawnPainting struct {
	EntityID           int32
	Title              string
	X, Y, Z, Direction int32
}
type SpawnExperienceOrb struct {
	EntityID int32
	X, Y, Z  int32
	Count    int16
}
type EntityVelocity struct {
	EntityID int32
	X, Y, Z  int16
}
type DestroyEntity struct {
	EntityIDs []int32
}
type CreateEntity struct {
	EntityID int32
}
type EntityRelativeMove struct {
	EntityID   int32
	DX, DY, DZ int8
}
type EntityLook struct {
	EntityID   int32
	Yaw, Pitch int8
}
type EntityLookRelativeMove struct {
	EntityID   int32
	DX, DY, DZ int8
	Yaw, Pitch int8
}
type EntityTeleport struct {
	EntityID   int32
	X, Y, Z    int32
	Yaw, Pitch int8
}
type EntityHeadLook struct {
	EntityID int32
	HeadYaw  byte
}
type EntityStatus struct {
	EntityID int32
	Status   EntityStatusType
}
type AttachEntity struct {
	EntityID  int32
	VehicleID int32
}
type SetEntityMetadata struct {
	EntityID int32
	Metadata []EntityMetadata
}
type EntityEffect struct {
	EntityID  int32
	EffectID  int8 // http://www.minecraftwiki.net/wiki/Potion_effect#Parameters
	Amplifier int8
	Duration  int16
}
type RemoveEntityEffect struct {
	EntityID int32
	EffectID int8 // http://www.minecraftwiki.net/wiki/Potion_effect#Parameters
}
type SetExperience struct {
	Percent float32
	Level   int16
	Total   int16
}
type ChunkData struct {
	X, Y                 int32
	IsGroundUpContinuous bool
	PrimaryBitMap        int16
	AddBitMap            int16
	ZlibData             Int32PrefixedBytes // compressed via zlib; see http://www.wiki.vg/Map_Format
}
type MultiBlockChange struct {
	ChunkX, ChunkY int32
	RecordCount    int16
	Data           Int32PrefixedBytes // see http://www.wiki.vg/Protocol#Multi_Block_Change_.280x34.29
}
type BlockChange struct {
	X        int32
	Y        int8
	Z        int32
	Type     int16
	Metadata int8
}
type BlockAction struct {
	X               int32
	Y               int16
	Z               int32
	InstrumentType  int8
	InstrumentPitch int8
	BlockID         int16
}
type BlockBreakAnimation struct {
	EntityID     int32
	X, Y, Z      int32
	DestroyStage int8
}
type MapChunkBulk struct {
	SkylightSent   bool
	CompressedData []byte // compressed
	Metadatas      []ChunkBulkMetadata
}
type Explosion struct {
	X, Y, Z                                           float64
	Radius                                            float32
	AffectedBlocks                                    []BlockPosition
	PlayerXVelocity, PlayerYVelocity, PlayerZVelocity float32
}
type Effect struct {
	EffectID         int32
	X                int32
	Y                int8
	Z                int32
	Data             int32 // depends on effect id
	NoVolumeDecrease bool
}
type NamedSoundEffect struct {
	Name    string
	X, Y, Z int32
	Volume  float32
	Pitch   int8
}
type ChangeGameState struct {
	State    GameState
	GameMode GameMode
}
type GlobalEntity struct {
	EntityID int32
	ID       int8  // = 1 for thunderbolt
	X, Y, Z  int32 // absolute
}
type OpenWindow struct {
	WindowID      int8
	InventoryType int8
	Title         string
	NumSlots      int8
}
type CloseWindow struct {
	WindowID int8
}
type ClickWindow struct {
	WindowID     int8
	Slot         int16
	MouseButton  MouseButton
	ActionNumber int16
	ShiftPressed bool
	ClickedItem  Slot
}
type SetSlot struct {
	WindowID int8
	Slot     int16
	Data     Slot
}
type SetWindowItems struct {
	WindowID int8
	Slots    []Slot
}
type UpdateWindowProperty struct {
	WindowID        int8
	Property, Value int16
}
type ConfirmTransaction struct {
	WindowID     int8
	ActionNumber int16
	Accepted     bool
}
type CreativeInventoryAction struct {
	Slot        int16
	ClickedItem Slot
}
type EnchantItem struct {
	WindowID    int8
	Enchantment int8
}
type UpdateSign struct {
	X                          int32
	Y                          int16
	Z                          int32
	Line1, Line2, Line3, Line4 string
}
type ItemData struct {
	Type int16
	ID   int16
	Text string // ascii string
}
type UpdateTileEntity struct {
	X       int32
	Y       int16
	Z       int32
	Action  byte
	NBTData []byte
}
type IncrementStatistic struct {
	ID     int32
	Amount int8
}
type PlayerListItem struct {
	Name   string
	Online bool
	Ping   int16
}
type PlayerAbilities struct {
	Flags        int8
	FlyingSpeed  int8
	WalkingSpeed int8
}
type TabComplete struct {
	Text string
}
type ClientSettings struct {
	Locale     string
	ViewDist   ViewDistance
	ChatFlags  int8
	Difficulty GameDifficulty
	ShowCape   bool
}
type ClientStatus struct {
	Payload int8 // 0 = initial spawn, 1 = respawn
}
type PluginMessage struct {
	Channel string
	Data    []byte
}
type EncryptionKeyResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}
type EncryptionKeyRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken []byte
}
type ServerListPing struct {
	Magic int8 // should always equal 1
}
type Disconnect struct {
	Reason string
}

type ScoreboardObjective struct {
	Name  string
	Value string
	Type  ScoreboardType
}

type UpdateScore struct {
	ItemName  string
	Type      ScoreType
	ScoreName string // only sent if not removing
	Value     int32  // only sent if not removing
}

type DisplayScoreboard struct {
	Position ScoreboardPosition
	Name     string
}

type Teams struct {
	ID           string
	Type         TeamType
	Name         string               // only for create or update
	Prefix       string               // only for create or update
	FriendlyFire TeamFriendlyFireType // only for create or update
	PlayersDelta []string             // only for create or add/remove players
}

///////////////////////////////////////////////////////

type BlockPosition struct {
	X, Y, Z int32
}

type ChunkBulkMetadata struct {
	ChunkX, ChunkY int32
	PrimaryBitmap  int16
	AddBitmap      int16 // unused?
}

// represents an item
type Slot struct {
	ID         int16
	Count      int8
	Damage     int16
	GzippedNBT []byte
}

// EmptySlot represents a special instance of the Slot type that indicates
// a slot with no items.
var EmptySlot = Slot{ID: -1}

func (s *Slot) IsEmpty() bool {
	return s.ID == -1
}

func (s *Slot) NewReader() (io.Reader, error) {
	return gzip.NewReader(bytes.NewReader(s.GzippedNBT))
}

// requires special parsing
// see http://www.wiki.vg/Entities#Entity_Metadata_Format
// clients require at least one piece of metadata (index 0, 1, or 8)
// indicates various entity states:
// bit index | bit mask | meaning
// 0         | 0x01     | Entity on Fire
// 1         | 0x02     | Entity crouched
// 2         | 0x04     | Entity riding
// 3         | 0x08     | Entity sprinting
// 4         | 0x10     | Eating/Drinking/Blocking/RightClickActions
//
// This represents one entry
type EntityMetadata struct {
	ID    EntityMetadataIndex
	Type  EntityMetadataType
	Value interface{}
}

type Int32PrefixedBytes []byte

//////////////////////////////////////////////////////
// currently for position parsing of EntityMetadataType 0x06
type Position struct {
	X, Y, Z int32
}

//////////////////////////////////////////////////////

type EntityMetadataIndex byte

const (
	EntityFlags         EntityMetadataIndex = 0
	EntityDrowning                          = 1
	EntityUnderPotionFX                     = 8
	EntityAnimalCounter                     = 12
	EntityState1                            = 16 // ie - creeper fuse, dragon hp, slime size etc.
	EntityState2                            = 17 // ie - enderman item metadata
	EntityState3                            = 18 // ie - enderman aggression
	EntityState4                            = 19 // ie - wolf unknown, minecart damage
)

//////////////////////////////////////////////////////

type EntityMetadataType byte

const (
	EntityMetadataByte EntityMetadataType = iota
	EntityMetadataShort
	EntityMetadataInt
	EntityMetadataFloat
	EntityMetadataString
	EntityMetadataSlot
	EntityMetadataPosition // x, y, z int
)

//////////////////////////////////////////////////////

type PacketType byte

///////////////////////////////////////////////////////

const (
	DefaultLevelType     string = "default"
	FlatLevelType               = "flat"
	LargeBiomesLevelType        = "largeBiomes"
)

type ViewDistance int8

const (
	FarViewDistance ViewDistance = iota
	NormalViewDistance
	ShortViewDistance
	TinyViewDistance
)

type GameState int8

const (
	GameStateInvalidBed GameState = iota
	GameStateBeginRain
	GameStateEndRain
	GameStateChangeGameMode
	GameStateEnterCredits
)

type EntityStatusType byte

const (
	EntityStatusHurt         EntityStatusType = 2
	EntityStatusDead                          = 3
	EntityStatusTaming                        = 6
	EntityStatusTamed                         = 7
	EntityStatusShakingWater                  = 8
	EntityStatusEating                        = 9
	EntityStatusEatingGrass                   = 10
)

type MobType int8

const (
	MobCreeper      MobType = 50
	MobSkeleton             = 51
	MobSpider               = 52
	MobGiantZombie          = 53
	MobZombie               = 54
	MobSlime                = 55
	MobGhast                = 56
	MobZombiePigman         = 57
	MobEnterman             = 58
	MobCaveSpider           = 59
	MobSilverFish           = 60
	MobBlaze                = 61
	MobMagmaCube            = 62
	MobEnderDragon          = 63
	MobWither               = 64
	MobBat                  = 65
	MobWitch                = 66
	MobPig                  = 90
	MobSheep                = 91
	MobCow                  = 92
	MobChicken              = 93
	MobSquid                = 94
	MobWolf                 = 95
	MobMooshroom            = 96
	MobSnowman              = 97
	MobOcelot               = 98
	MobIronGolem            = 99
	MobVillager             = 120
)

type EntityType int8

const (
	EntityBoat             EntityType = 1
	EntityMinecart                    = 10
	EntityMinecartStorage             = 11
	EntityMinecartPowered             = 12
	EntityActiveTNT                   = 50
	EntityEnderCrystal                = 51
	EntityArrow                       = 60
	EntitySnowball                    = 61
	EntityEgg                         = 62
	EntityEnderpearl                  = 65
	EntityWitherSkull                 = 66
	EntityFallingObject               = 70
	EntityEyeOfEnder                  = 72
	EntityThrownPotion                = 73
	EntityFallingDragonEgg            = 74
	EntityThrownExpBottle             = 75
	EntityFishingFloat                = 90
)

type ActionType int8

const (
	ActionCrouch ActionType = iota
	ActionUncrouch
	ActionLeaveBed
	ActionStartSprinting
	ActionStopSprinting
)

type AnimationType int8

const (
	AnimationNone     AnimationType = 0
	AnimationSwingArm               = 1 // only clients send this
	AnimationDamage                 = 2 // only server sends this
	AnimationLeaveBed               = 3
	AnimationEatFood                = 5
	AnimationUnknown                = 102
	AnimationCrouch                 = 104
	AnimationUncrouch               = 105
)

type Face int8

const (
	FaceYNeg Face = iota
	FaceYPos
	FaceZNeg
	FaceZPos
	FaceXNeg
	FaceXPos
)

type PlayerDiggingStatus int8

const (
	PlayerStartedDigging           PlayerDiggingStatus = iota // can also open doors
	PlayerCancelledDigging                                    // not sent by client
	PlayerFinishedDigging                                     // client sends when it thinks it's finished
	PlayerCheckBlock                                          // not sent by client
	PlayerDropItem                                            // zero all other fields for PlayerDigging
	PlayerShootArrowOrFinishEating                            // zero all other fields for PlayerDigging, except Face = 255
)

type MouseButton int8

const (
	ButtonLeftMouse MouseButton = iota
	ButtonRightMouse
	ButtonShift
	ButtonMiddleMouse
)

type GameDifficulty uint8

const (
	GameDifficultyPeaceful GameDifficulty = iota
	GameDifficultyEasy
	GameDifficultyNormal
	GameDifficultyHard
)

///////////////////////////////////////////////////////

type GameDimension int8

func (d *GameDimension) IsNether() bool    { return *d == GameDimensionNether }
func (d *GameDimension) IsOverworld() bool { return *d == GameDimensionOverworld }
func (d *GameDimension) IsEnd() bool       { return *d == GameDimensionEnd }

const (
	GameDimensionNether GameDimension = iota - 1
	GameDimensionOverworld
	GameDimensionEnd
)

///////////////////////////////////////////////////////

type GameMode int8

func (m *GameMode) IsSurvival() bool  { return *m&GameModeFlag == GameModeSurvival }
func (m *GameMode) IsCreative() bool  { return *m&GameModeFlag == GameModeCreative }
func (m *GameMode) IsAdventure() bool { return *m&GameModeFlag == GameModeAdventure }
func (m *GameMode) IsHardcore() bool  { return *m&GameModeHardcoreFlag > 0 }
func (m *GameMode) Name() string {
	switch *m {
	case GameModeSurvival:
		return "Survival"
	case GameModeCreative:
		return "Creative"
	case GameModeAdventure:
		return "Adventure"
	}
	return "Unknown"
}

const (
	GameModeSurvival GameMode = iota
	GameModeCreative
	GameModeAdventure
	GameModeFlag         = 0x3
	GameModeHardcoreFlag = 0x8
)

///////////////////////////////////////////////////////

type ScoreboardType byte

const (
	ScoreboardTypeCreate ScoreboardType = iota
	ScoreboardTypeDelete
	ScoreboardTypeUpdate
)

type ScoreType byte

const (
	ScoreTypeCreateOrUpdate ScoreType = iota
	ScoreTypeDelete
)

type ScoreboardPosition byte

const (
	ScoreboardPositionList ScoreboardPosition = iota
	ScoreboardPositionSidebar
	ScoreboardPositionBelowName
)

///////////////////////////////////////////////////////

type TeamType byte

const (
	TeamCreate TeamType = iota
	TeamDelete
	TeamUpdate
	TeamPlayerAdd
	TeamPlayerDelete
)

type TeamFriendlyFireType byte

const (
	TeamFriendlyFireOff TeamFriendlyFireType = iota
	TeamFriendlyFireOn
	TeamFriendlyFireShowFriendlyInvisible = 3
)
