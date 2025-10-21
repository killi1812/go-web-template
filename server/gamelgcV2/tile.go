package gamelgcV2

type tile struct {
	TileId   int          `json:"tileId"`   // id of the tile
	Category tileCategory `json:"category"` // Category describes category of tiles (suited, honors, bonus)

	// Suite represents a suite of the tile (bamboo, dot, wind, dragon ...)
	//
	// ---!!! Not to be confused with suited category !!!---
	Suite   tileSuite `json:"suite"`
	Value   int       `json:"value"`   // Value represents a value of the tile
	Visible bool      `json:"visible"` // Visible true if the tile is Visible by all players
	// NOTE: change to method .IsBonus()
	IsBonus bool `json:"isBonus"` // Bonus describes if the tile is the bonus tile for extra points
}

func (t tile) Equal(other tile) bool {
	return t.Category == other.Category && t.Suite == other.Suite && t.Value == other.Value
}

type (
	tileCategory string // one of suited, bonus, honor
	tileSuite    string // one of bamboo, dot, character, wind, dragon, season, flower
	valueWind    = int  // values of wind tiles
	valueDragon  = int  // values of dragon tiles
	valueFlower  = int  // values of flowers tiles
	valueSeason  = int  // values of seasons tiles
)

// Main tile categoryes
const (
	_CATEGORY_SUITED tileCategory = "suited"
	_CATEGORY_HONOR  tileCategory = "honor"
	_CATEGORY_BONUS  tileCategory = "bonus"
)

// Suited Tiles
const (
	_SUITE_DOT       tileSuite = "dot"
	_SUITE_BAMBOO    tileSuite = "bamboo"
	_SUITE_CHARACTER tileSuite = "character"
)

// Honor tiles
const (
	_SUITE_WIND   tileSuite = "wind"
	_SUITE_DRAGON tileSuite = "dragon"
)

// Bonus tiles
const (
	_SUITE_FLOWER tileSuite = "flower"
	_SUITE_SEASON tileSuite = "season"
)

const (
	// Honor Tiles - Winds

	_VALUE_WIND_EAST  valueWind = 11 + iota
	_VALUE_WIND_SOUTH valueWind = 11 + iota
	_VALUE_WIND_WEST  valueWind = 11 + iota
	_VALUE_WIND_NORTH valueWind = 11 + iota

	// Honor Tiles - Dragons

	_VALUE_DRAGON_WHITE valueDragon = 11 + iota
	_VALUE_DRAGON_GREEN valueDragon = 11 + iota
	_VALUE_DRAGON_RED   valueDragon = 11 + iota

	// Bonus tiles - Flowers

	_VALUE_FLOWER_PLUM_BLOSSOM  valueFlower = 11 + iota
	_VALUE_FLOWER_ORCHID        valueFlower = 11 + iota
	_VALUE_FLOWER_CHRYSANTHEMUM valueFlower = 11 + iota
	_VALUE_FLOWER_BAMBOO        valueFlower = 11 + iota

	// Bonus tiles - Flowers

	_VALUE_SEASON_SPRING valueSeason = 11 + iota
	_VALUE_SEASON_SUMMER valueSeason = 11 + iota
	_VALUE_SEASON_AUTUMN valueSeason = 11 + iota
	_VALUE_SEASON_WINTER valueSeason = 11 + iota
)

var _ALL_SUITED_SUITS = []tileSuite{_SUITE_BAMBOO, _SUITE_DOT, _SUITE_CHARACTER}
var _ALL_HONORS_SUITS = []tileSuite{_SUITE_DRAGON, _SUITE_WIND}
var _ALL_BONUS_SUITS = []tileSuite{_SUITE_FLOWER, _SUITE_SEASON}

var _ALL_DRAGONS_VALUES = []valueDragon{_VALUE_DRAGON_GREEN, _VALUE_DRAGON_RED, _VALUE_DRAGON_WHITE}
var _ALL_WINDS_VALUES = []valueDragon{_VALUE_WIND_EAST, _VALUE_WIND_NORTH, _VALUE_WIND_WEST, _VALUE_WIND_SOUTH}
var _ALL_HONORS_VALUES = append(_ALL_DRAGONS_VALUES, _ALL_WINDS_VALUES...)

var _ALL_SEASONS_VALUES = []valueSeason{_VALUE_SEASON_AUTUMN, _VALUE_SEASON_SPRING, _VALUE_SEASON_SUMMER, _VALUE_SEASON_WINTER}
var _ALL_FLOWERS_VALUES = []valueFlower{_VALUE_FLOWER_BAMBOO, _VALUE_FLOWER_CHRYSANTHEMUM, _VALUE_FLOWER_ORCHID, _VALUE_FLOWER_PLUM_BLOSSOM}
var _ALL_BONUS_VALUES = append(_ALL_SEASONS_VALUES, _ALL_FLOWERS_VALUES...)
