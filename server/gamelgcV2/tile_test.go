package gamelgcV2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTileEqual(t *testing.T) {
	mainT := tile{
		TileId:   1,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_BAMBOO,
		Value:    1,
	}
	eqT := tile{
		TileId:   2,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_BAMBOO,
		Value:    1,
	}
	honorT := tile{
		TileId:   3,
		Category: _CATEGORY_HONOR,
		Suite:    _SUITE_DRAGON,
		Value:    15,
	}
	b2T := tile{
		TileId:   4,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_BAMBOO,
		Value:    2,
	}
	c2T := tile{
		TileId:   5,
		Category: _CATEGORY_SUITED,
		Suite:    _SUITE_CHARACTER,
		Value:    2,
	}

	assert.True(t, mainT.Equal(eqT), "Tiles should be equal")
	assert.False(t, mainT.Equal(honorT), "Tiles should not be equal different category")
	assert.False(t, mainT.Equal(b2T), "Tiles should not be equal different values")
	assert.False(t, mainT.Equal(c2T), "Tiles should not be equal different suite")
}
