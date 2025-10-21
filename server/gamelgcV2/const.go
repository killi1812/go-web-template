package gamelgcV2

import "time"

// _HAND_SIZE is default hand size can be more but not less
const _HAND_SIZE = 13

// _WALL_SIZE is max size of the wall
const _WALL_SIZE = 144

// _PLAYER_COUNT is count of the players in the game
const _PLAYER_COUNT = 4

// _RESPONSE_TIME is a default time for player to response
const _RESPONSE_TIME = 3 * time.Second

type gameAction string

// _ACTION_PUNG_TILE
//
// desc: take a face down tile
const _ACTION_DRAW_TILE gameAction = "draw_tile"

// _ACTION_DISCARD_TILE
//
// desc: drop an extra tile
const _ACTION_DISCARD_TILE gameAction = "discard_tile"

// _ACTION_PASS_TURN
//
// desc: pass the turn on to the next player
const _ACTION_PASS_TURN gameAction = "pass_turn"

// _ACTION_DECLARE_RICHI
//
// desc: declare that you have richi
const _ACTION_DECLARE_RICHI gameAction = "pass_turn"

// _ACTION_SEIZE
//
// desc: take a tile from discard pile, order of power is from least to most:pung, chow, kung
const _ACTION_SEIZE gameAction = "seize"

type seizeType int

const (
	// _SEIZE_CHOW
	//
	// desc: take a tile from discard to complete 3 tiles of the same suite in a sequence
	_SEIZE_CHOW seizeType = 1 + iota

	// _SEIZE_PUNG
	//
	// desc: take a tile from discard to complete 3 identical tiles
	_SEIZE_PUNG seizeType = 1 + iota

	// _SEIZE_KUNG
	//
	// desc: take a tile from discard to complete 4 identical tiles
	_SEIZE_KUNG seizeType = 1 + iota
)
