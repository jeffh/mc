package protocol

import (
	"fmt"
)

// protocol version supported
const Version = 49

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
	basePacketMapper.Define(0x28, EntityMetadata{})
	basePacketMapper.Define(0x29, EntityEffect{})
	basePacketMapper.Define(0x2A, RemoveEntityEffect{})
	basePacketMapper.Define(0x2B, SetExperience{})
	basePacketMapper.Define(0x33, ChunkData{})
	basePacketMapper.Define(0x34, MultiBlockChange{})
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
	Message string
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
type PlayerPositionLookForServer struct {
	X, Y, Stance, Z float64
	Yaw, Pitch      float32
	IsOnGround      bool
}
type PlayerPositionLookForClient struct {
	X, Stance, Y, Z float64
	Yaw, Pitch      float32
	IsOnGround      bool
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
	Metadata    EntityMetadata
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
	Yaw, Pitch, HeadYaw             int8
	ZVelocity, XVelocity, YVelocity int16
	Metadata                        EntityMetadata
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
	HeadYaw  int8
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
	Metadata EntityMetadata
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
	Level   int8
	Total   int8
}
type ChunkData struct {
	X, Y                 int32
	IsGroundUpContinuous bool
	PrimaryBitMap        uint16
	AddBitMap            uint16
	CompressedData       []uint8 // compressed via zlib; see http://www.wiki.vg/Map_Format
}
type MultiBlockChange struct {
	ChunkX, ChunkY int32
	RecordCount    int16
	Data           []int32 // see http://www.wiki.vg/Protocol#Multi_Block_Change_.280x34.29
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
	Count          int16
	CompressedData []byte // compressed
	Metadata       []ChunkBulkMetadata
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
	Count    int16
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
	Action  int8
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
	ID     int16
	Count  int8
	Damage int16
	Data   []byte
}

// requires special parsing
// see http://www.wiki.vg/Entities#Entity_Metadata_Format
// clients require at least one piece of metadata (index 0, 1, or 8)
type EntityMetadata struct {
	// indicates various entity states:
	// bit index | bit mask | meaning
	// 0         | 0x01     | Entity on Fire
	// 1         | 0x02     | Entity crouched
	// 2         | 0x04     | Entity riding
	// 3         | 0x08     | Entity sprinting
	// 4         | 0x10     | Eating/Drinking/Blocking/RightClickActions
	Flags        int8  // index 0
	DrownCounter int8  // index 1: starts at 300 -> -19
	PotionEffect int32 // index 8: 0x00RRGGBB or 0 if no effects
	Animals      int32 // index 12: -23999 = baby -> 0 = normal <- 6000 = parent
}

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
	InvalidBedState GameState = iota
	BeginRainState
	EndRainState
	ChangeGameModeState
	EnterCredits
)

type EntityStatusType int8

const (
	HurtStatusType         EntityStatusType = 2
	DeadStatusType                          = 3
	TamingStatusType                        = 6
	TamedStatusType                         = 7
	ShakingWaterStatusType                  = 8
	EatingStatusType                        = 9
	EatingGrassStatusType                   = 10
)

type MobType int8

const (
	CreeperMobType      MobType = 50
	SkeletonMobType             = 51
	SpiderMobType               = 52
	GiantZombieMobType          = 53
	ZombieMobType               = 54
	SlimeMobType                = 55
	GhastMobType                = 56
	ZombiePigmanMobType         = 57
	EntermanMobType             = 58
	CaveSpiderMobType           = 59
	SilverFishMobType           = 60
	BlazeMobType                = 61
	MagmaCubeMobType            = 62
	EnderDragonMobType          = 63
	WitherMobType               = 64
	BatMobType                  = 65
	WitchMobType                = 66
	PigMobType                  = 90
	SheepMobType                = 91
	CowMobType                  = 92
	ChickenMobType              = 93
	SquidMobType                = 94
	WolfMobType                 = 95
	MooshroomMobType            = 96
	SnowmanMobType              = 97
	OcelotMobType               = 98
	IronGolemMobType            = 99
	VillagerMobType             = 120
)

type EntityType int8

const (
	BoatEntityType             EntityType = 1
	MinecartEntityType                    = 10
	MinecartStorageEntityType             = 11
	MinecartPoweredEntityType             = 12
	ActiveTNTEntityType                   = 50
	EnderCrystalEntityType                = 51
	ArrowEntityType                       = 60
	SnowballEntityType                    = 61
	EggEntityType                         = 62
	EnderpearlEntityType                  = 65
	WitherSkullEntityType                 = 66
	FallingObjectEntityType               = 70
	EyeOfEnderEntityType                  = 72
	ThrownPotionEntityType                = 73
	FallingDragonEggEntityType            = 74
	ThrownExpBottleEntityType             = 75
	FishingFloatEntityType                = 90
)

type ActionType int8

const (
	CrouchAction ActionType = iota
	UncrouchAction
	LeaveBedAction
	StartSprintingAction
	StopSprintingAction
)

type AnimationType int8

const (
	NoAnimation       AnimationType = 0
	SwingArmAnimation               = 1 // only clients send this
	DamageAnimation                 = 2 // only server sends this
	LeaveBedAnimation               = 3
	EatFoodAnimation                = 5
	UnknownAnimation                = 102
	CrouchAnimation                 = 104
	UncrouchAnimation               = 105
)

type Face int8

const (
	YNeg Face = iota
	YPos
	ZNeg
	ZPos
	XNeg
	XPos
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
	LeftMouseButton MouseButton = iota
	RightMouseButton
	ShiftButton
	MiddleMouseButton
)

type GameDifficulty uint8

const (
	PeacefulDifficulty GameDifficulty = iota
	EasyDifficulty
	NormalDifficulty
	HardDifficulty
)

///////////////////////////////////////////////////////

type GameDimension int8

func (d *GameDimension) IsNether() bool    { return *d == NetherDimension }
func (d *GameDimension) IsOverworld() bool { return *d == OverworldDimension }
func (d *GameDimension) IsEnd() bool       { return *d == EndDimension }

const (
	NetherDimension GameDimension = iota - 1
	OverworldDimension
	EndDimension
)

///////////////////////////////////////////////////////

type GameMode int8

func (m *GameMode) IsSurvival() bool  { return *m&GameModeFlag == SurvivalMode }
func (m *GameMode) IsCreative() bool  { return *m&GameModeFlag == CreativeMode }
func (m *GameMode) IsAdventure() bool { return *m&GameModeFlag == AdventureMode }
func (m *GameMode) IsHardcore() bool  { return *m&HardcoreModeFlag > 0 }
func (m *GameMode) Name() string {
	switch *m {
	case SurvivalMode:
		return "Survival"
	case CreativeMode:
		return "Creative"
	case AdventureMode:
		return "Adventure"
	}
	return "Unknown"
}

const (
	SurvivalMode GameMode = iota
	CreativeMode
	AdventureMode
	GameModeFlag     = 0x3
	HardcoreModeFlag = 0x8
)
