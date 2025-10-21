package gamelgcV2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscard(t *testing.T) {
	// 1. Setup initial player state with a known hand
	playerState := &PlayerState{
		MainHand: hand{
			{TileId: 1, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 1},
			{TileId: 2, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 2},
			{TileId: 3, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 3},
		},
	}
	initialHandSize := len(playerState.MainHand)
	tileToDiscardId := 2

	// 2. Test a valid discard
	discardedTile := playerState.discard(tileToDiscardId)

	// Assert that the correct tile was returned
	assert.NotNil(t, discardedTile, "Discarded tile should not be nil for a valid discard")
	assert.Equal(t, tileToDiscardId, discardedTile.TileId, "The correct tile should be discarded")

	// Assert that the hand size has decreased by one
	assert.Len(t, playerState.MainHand, initialHandSize-1, "Hand size should decrease by one after a valid discard")

	// Assert that the discarded tile is no longer in the hand
	for _, tile := range playerState.MainHand {
		assert.NotEqual(t, tileToDiscardId, tile.TileId, "Discarded tile should not be in the hand anymore")
	}

	// 3. Test discarding a tile that is not in the hand
	invalidTileId := 99
	nonDiscardedTile := playerState.discard(invalidTileId)

	// Assert that nil is returned for an invalid discard
	assert.Nil(t, nonDiscardedTile, "Should return nil when trying to discard a tile that is not in hand")

	// Assert that the hand size remains unchanged
	assert.Len(t, playerState.MainHand, initialHandSize-1, "Hand size should not change after an invalid discard attempt")
}

func TestHasPungWith(t *testing.T) {
	pungTile := tile{
		TileId:   11,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_BAMBOO,
		Value:    1,
	}
	noPungTile := tile{
		TileId:   10,
		Category: _CATEGORY_HONOR,
		Suite:    _SUITE_DRAGON,
		Value:    _VALUE_DRAGON_GREEN,
	}

	var thand hand = []tile{
		{
			TileId:   1,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    1,
		},
		{
			TileId:   2,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    1,
		},
		{
			TileId:   3,
			Category: _CATEGORY_HONOR,
			Suite:    _SUITE_DRAGON,
			Value:    _VALUE_DRAGON_GREEN,
		},
		{
			TileId:   4,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    2,
		},
		{
			TileId:   5,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_CHARACTER,
			Value:    2,
		},
	}

	pungGood := thand.HasPungWith(pungTile)
	assert.Len(t, pungGood, 2, "Return len should be 2")
	for _, tile := range pungGood {
		assert.NotEqual(t, pungTile.TileId, tile.TileId, "Tiles should not have sae id as main tile")
	}

	pungBad := thand.HasPungWith(noPungTile)
	assert.Nil(t, pungBad, "Return value should be nil")
}

func TestHasKungWith(t *testing.T) {
	pungTile := tile{
		TileId:   11,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_BAMBOO,
		Value:    1,
	}
	noPungTile := tile{
		TileId:   10,
		Category: _CATEGORY_HONOR,
		Suite:    _SUITE_DRAGON,
		Value:    _VALUE_DRAGON_GREEN,
	}

	var thand hand = []tile{
		{
			TileId:   1,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    1,
		},
		{
			TileId:   2,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    1,
		},
		{
			TileId:   3,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    1,
		},

		{
			TileId:   4,
			Category: _CATEGORY_HONOR,
			Suite:    _SUITE_DRAGON,
			Value:    _VALUE_DRAGON_GREEN,
		},
		{
			TileId:   5,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_BAMBOO,
			Value:    2,
		},
		{
			TileId:   6,
			Category: _CATEGORY_SUITED,
			Suite:    _SUITE_CHARACTER,
			Value:    2,
		},
	}

	kungGood := thand.HasKungWith(pungTile)
	assert.Len(t, kungGood, 3, "Return len should be 3")
	for _, tile := range kungGood {
		assert.NotEqual(t, pungTile.TileId, tile.TileId, "Tiles should not have sae id as main tile")
	}

	// TODO: move this to test pung but it ok
	pung := thand.HasPungWith(noPungTile)
	assert.Nil(t, pung, "Return value should be nil")

	kungBad := thand.HasKungWith(noPungTile)
	assert.Nil(t, kungBad, "Return value should be nil")
}

func TestHasChowWith(t *testing.T) {
	// Define some tiles for testing
	bamboo3 := tile{TileId: 1, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 3}
	bamboo4 := tile{TileId: 2, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 4}
	bamboo5 := tile{TileId: 3, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 5}
	bamboo6 := tile{TileId: 4, Category: _CATEGORY_SUITED, Suite: _SUITE_BAMBOO, Value: 6}
	dot4 := tile{TileId: 5, Category: _CATEGORY_SUITED, Suite: _SUITE_DOT, Value: 4}
	eastWind := tile{TileId: 6, Category: _CATEGORY_HONOR, Suite: _SUITE_WIND, Value: _VALUE_WIND_EAST}

	tests := []struct {
		name         string
		hand         hand
		incomingTile tile
		want         []tile
	}{
		{
			name:         "Success Case: Incoming tile is the lowest (3)",
			hand:         hand{bamboo4, bamboo5},
			incomingTile: bamboo3,
			want:         []tile{bamboo4, bamboo5},
		},
		{
			name:         "Success Case: Incoming tile is in the middle (4)",
			hand:         hand{bamboo3, bamboo5},
			incomingTile: bamboo4,
			want:         []tile{bamboo3, bamboo5},
		},
		{
			name:         "Success Case: Incoming tile is the highest (5)",
			hand:         hand{bamboo3, bamboo4},
			incomingTile: bamboo5,
			want:         []tile{bamboo3, bamboo4},
		},
		{
			name:         "Failure Case: No chow possible",
			hand:         hand{bamboo3, bamboo6},
			incomingTile: bamboo4,
			want:         nil,
		},
		{
			name:         "Failure Case: Correct values but wrong suit",
			hand:         hand{bamboo3, dot4},
			incomingTile: bamboo5,
			want:         nil,
		},
		{
			name:         "Failure Case: Incoming tile is not a suited tile (Honor)",
			hand:         hand{bamboo3, bamboo4},
			incomingTile: eastWind,
			want:         nil,
		},
		{
			name:         "Failure Case: Hand is empty",
			hand:         hand{},
			incomingTile: bamboo4,
			want:         nil,
		},
		{
			name:         "Edge Case: Sequence at the beginning of a suit (1-2-3)",
			hand:         hand{{Value: 1, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}, {Value: 2, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}},
			incomingTile: tile{Value: 3, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED},
			want:         []tile{{Value: 1, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}, {Value: 2, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}},
		},
		{
			name:         "Edge Case: Sequence at the end of a suit (7-8-9)",
			hand:         hand{{Value: 8, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}, {Value: 9, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}},
			incomingTile: tile{Value: 7, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED},
			want:         []tile{{Value: 8, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}, {Value: 9, Suite: _SUITE_DOT, Category: _CATEGORY_SUITED}},
		},
		{
			name:         "Failure case: Hand has duplicates, but not a valid chow",
			hand:         hand{bamboo3, bamboo3, bamboo5},
			incomingTile: bamboo4,
			want:         []tile{bamboo3, bamboo5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hand.HasChowWith(tt.incomingTile)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				// Use ElementsMatch to compare slices without worrying about order.
				assert.ElementsMatch(t, tt.want, got)
			}
		})
	}
}
