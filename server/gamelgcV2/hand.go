package gamelgcV2

import (
	"slices"
	"template/model"

	"go.uber.org/zap"
)

type PlayerState struct {
	user *model.User
	wind valueWind // Wind represents Player's Wind

	MainHand     hand   `json:"mainHand"`     // MainHand represents a players hand
	RevealedHand []tile `json:"revealedHand"` // RevealedHand represents players revealed tiles
	BonusTiles   []tile `json:"bonusTiles"`   // Bonus represents players bonus tiles
}

// CheckSeize will check if the seize is possible
func (h *PlayerState) CheckSeize(tle *tile) ([]tile, seizeType) {
	var seize []tile

	seize = h.MainHand.HasPungWith(*tle)
	if seize != nil {
		return seize, _SEIZE_PUNG
	}

	seize = h.MainHand.HasChowWith(*tle)
	if seize != nil {
		return seize, _SEIZE_CHOW
	}

	seize = h.MainHand.HasKungWith(*tle)
	if seize != nil {
		return seize, _SEIZE_CHOW
	}
	return nil, 0
}

// discard removes a tile from the hand of the player and returns it
func (h *PlayerState) discard(tileId int) *tile {

	findFunc := func(t tile) bool {
		return t.TileId == tileId
	}

	// find discarded tile
	index := slices.IndexFunc(h.MainHand, findFunc)
	if index == -1 {
		zap.S().Errorf("Tile not in hand,tileId: %d", tileId)
		return nil
	}

	discardedTile := h.MainHand[index]

	// update the hand
	h.MainHand = append(h.MainHand[:index], h.MainHand[index+1:]...)

	return &discardedTile
}

type hand []tile

// HasChowWith will return the two tiles that complete a chow with a given tile.
// It checks for all three possible sequence positions (e.g., 1-2-3, 2-3-4, 3-4-5).
// If no chow can be completed, it returns nil.
//
// ---!!! IMPORTANT !!!---
//
// This function does NOT check if the player is allowed to claim a chow (e.g., turn order).
// It only checks the mathematical possibility based on the hand's contents.
func (h hand) HasChowWith(tle tile) []tile {
	// Chows are only possible with suited tiles (bamboo, dots, characters).
	if tle.Category != _CATEGORY_SUITED {
		return nil
	}

	// Create a map for quick lookups of tiles in the hand that match the suit.
	suitedTiles := make(map[int]tile)
	for _, t := range h {
		if t.Suite == tle.Suite {
			suitedTiles[t.Value] = t
		}
	}

	// Case 1: The incoming tile is the lowest in the sequence (e.g., hand has 4, 5 and incoming is 3)
	if tile1, ok1 := suitedTiles[tle.Value+1]; ok1 {
		if tile2, ok2 := suitedTiles[tle.Value+2]; ok2 {
			return []tile{tile1, tile2}
		}
	}

	// Case 2: The incoming tile is in the middle of the sequence (e.g., hand has 3, 5 and incoming is 4)
	if tile1, ok1 := suitedTiles[tle.Value-1]; ok1 {
		if tile2, ok2 := suitedTiles[tle.Value+1]; ok2 {
			return []tile{tile1, tile2}
		}
	}

	// Case 3: The incoming tile is the highest in the sequence (e.g., hand has 3, 4 and incoming is 5)
	if tile1, ok1 := suitedTiles[tle.Value-2]; ok1 {
		if tile2, ok2 := suitedTiles[tle.Value-1]; ok2 {
			return []tile{tile1, tile2}
		}
	}

	// If no sequence was found, return nil.
	return nil
}

// HasKungWith will return tiles that complete kung with a given tile
// if no kung can be completed returned value will be nil
func (h hand) HasKungWith(tle tile) []tile {
	tiles := make([]tile, 0, 3)
	for _, t := range h {
		if t.Equal(tle) {
			tiles = append(tiles, t)
		}
	}
	if len(tiles) == 3 {
		return tiles
	}
	return nil
}

// HasPungWith will return tiles that complete pung with given tile
// if no pung can be completed returned value will be nil
func (h hand) HasPungWith(tle tile) []tile {
	tiles := make([]tile, 0, 2)
	for _, t := range h {
		if t.Equal(tle) {
			tiles = append(tiles, t)
		}
	}
	if len(tiles) == 2 {
		return tiles
	}
	return nil
}
