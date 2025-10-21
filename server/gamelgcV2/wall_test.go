package gamelgcV2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWall(t *testing.T) {
	// Call the function to create a new wall
	wall := newWall()

	// Assert that the wall has the correct number of tiles (144)
	assert.Len(t, wall, _WALL_SIZE, "A new wall should have 144 tiles")
}

func TestDraw(t *testing.T) {
	wall := wall{
		MainWall: []tile{
			{
				TileId: 1,
			},
			{
				TileId: 2,
			},
			{
				TileId: 3,
			},
			{
				TileId: 4,
			},
			{
				TileId: 5,
			},
			{
				TileId: 6,
			},
		},
		BackWall:  make([]tile, 10),
		BonusTile: tile{},
	}
	mainWallSize := len(wall.MainWall)

	// Test drawing a single tile
	drawnTiles := wall.draw()
	assert.Len(t, drawnTiles, 1, "Should draw one tile by default")
	assert.Equal(t, 1, drawnTiles[0].TileId, "First tile is drawn id should match")
	assert.Len(t, wall.MainWall, mainWallSize-1, "Wall should have one less tile after drawing one")

	// Test drawing another
	drawnTiles = wall.draw()
	assert.Len(t, drawnTiles, 1, "Should draw one tile by default")
	assert.Equal(t, 2, drawnTiles[0].TileId, "Second tile is drawn id should match")
	assert.Len(t, wall.MainWall, mainWallSize-2, "Wall should have one less tile after drawing one")

	// Test drawing multiple tiles
	drawnTiles = wall.draw(4)
	assert.Len(t, drawnTiles, 4, "Should draw the specified number of tiles")
	for i := range 4 {
		assert.Equal(t, i+3, drawnTiles[i].TileId, fmt.Sprintf("%d. tile drawn id should match", i+1))
	}
	assert.Len(t, wall.MainWall, mainWallSize-6, "Wall should have 5 less tiles in total")

	// Test drawing more tiles than are available
	drawnTiles = wall.draw(200)
	assert.Empty(t, drawnTiles, "Should not return any tiles when drawing more than available")

	// Test drawing zero tiles
	drawnTiles = wall.draw(0)
	assert.Empty(t, drawnTiles, "Should not return any tiles when drawing 0 tiles")

	// Test drawing negative number of tiles
	drawnTiles = wall.draw(-2)
	assert.Empty(t, drawnTiles, "Should not return any tiles when drawing negative number of tiles")
}

func TestDrawBack(t *testing.T) {
	wall := wall{
		BackWall: []tile{
			{
				TileId: 1,
			},
			{
				TileId: 6,
			},
		},
	}
	wallLen := len(wall.BackWall)

	// Test drawing from the back wall
	drawnTile := wall.drawBack()
	assert.NotEmpty(t, drawnTile, "Should draw one tile from the back wall")
	assert.Equal(t, 6, drawnTile.TileId, "First tile is drawn id should match")
	assert.Len(t, wall.BackWall, wallLen-1, "Back wall should have one less tile")

	wall.drawBack() // draw one more available tile
	drawnTile = wall.drawBack()
	assert.Empty(t, drawnTile, "Should not return any tiles when drawing more than available")

}

func TestShuffle(t *testing.T) {
	// Create a new, unshuffled wall
	wall := newWall()
	originalOrder := make([]tile, len(wall))
	copy(originalOrder, wall)

	// Shuffle the wall
	shuffledWall := shuffe(wall)

	// Assert that the length of the wall remains the same
	assert.Len(t, shuffledWall, _WALL_SIZE, "Shuffled wall should have the same number of tiles")

	// Assert that the order of tiles has changed
	assert.NotEqual(t, originalOrder, shuffledWall, "The order of tiles should be different after shuffling")
}

func TestBreakWall(t *testing.T) {
	// Create a new, unshuffled wall for a predictable test
	wall := newWall()

	// Break the wall
	brokenWall := breakWall(wall)

	// Assert that the total number of tiles is still correct
	assert.Equal(t, _WALL_SIZE, len(brokenWall.MainWall)+len(brokenWall.BackWall)+1, "Total number of tiles should remain the same after breaking the wall")

	// Assert that the BonusTile is not empty
	assert.NotZero(t, brokenWall.BonusTile.Suite, "Bonus tile should be set")
}
