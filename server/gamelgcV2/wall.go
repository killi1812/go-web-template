package gamelgcV2

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

// wall represents a wall of tile
type wall struct {
	MainWall  []tile `json:"mainWall"`  // Main represents a main drawing wall
	BonusTile tile   `json:"bonusTile"` // BonusTile is a revield tile
	BackWall  []tile `json:"backWall"`  // BackWall represents back of the wall after BonusTile
}

// draw draws tiles from the wall if the count is omitted then it draw 1 tile
func (w *wall) draw(count ...int) []tile {
	var cnt int
	if len(count) == 0 {
		cnt = 1
	} else {
		cnt = count[0]
	}

	if cnt <= 0 {
		zap.S().DPanicf("Trying to draw more then the len of the wall, count: %d. wall len %d", cnt, len(w.MainWall))
		return []tile{}
	}
	if cnt > len(w.MainWall) {
		zap.S().DPanicf("Trying to draw more then the len of the wall, count: %d. wall len %d", cnt, len(w.MainWall))
		return []tile{}
	}
	drawnTiles := w.MainWall[:cnt]
	w.MainWall = w.MainWall[cnt:]
	return drawnTiles
}

// draw draws one tile from back of the wall
func (w *wall) drawBack() tile {
	cnt := len(w.BackWall) - 1

	if cnt <= 0 {
		zap.S().DPanicf("Wall to small ", cnt, len(w.BackWall))
		return tile{}
	}
	drawnTiles := w.BackWall[cnt:]
	w.BackWall = w.BackWall[:cnt]
	return drawnTiles[0]
}

// TODO: gheck if the game logic is right
func breakWall(tiles []tile) wall {

	breakTileIndex := 7
	if breakTileIndex%2 == 1 {
		breakTileIndex++
	}

	breakTileIndex = (_WALL_SIZE - 1) - breakTileIndex

	var wall wall
	zap.S().Debugf("breaking the wall on index: %d", breakTileIndex)
	wall.BonusTile = tiles[breakTileIndex]
	wall.BackWall = tiles[breakTileIndex+1:]
	wall.MainWall = tiles[:breakTileIndex]
	return wall
}

// shuffle randomizes the order of the tiles in the wall.
// It uses the modern Fisher-Yates shuffle algorithm.
func shuffe(wall []tile) []tile {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform the shuffle
	r.Shuffle(len(wall), func(i, j int) {
		wall[i], wall[j] = wall[j], wall[i]
	})
	return wall
}

// newWall creates a standard 144-tile mahjong wall.
// This function acts as the "seed" for the game.
func newWall() []tile {
	tiles := make([]tile, 0, _WALL_SIZE)
	idCounter := 0

	// 1. Generate Suited Tiles (Dots, Bamboo, Characters)
	// There are 3 suits, with numbers 1-9. Each tile appears 4 times.
	for _, suit := range _ALL_SUITED_SUITS {
		for i := range 9 {
			for range 4 {
				tiles = append(tiles, tile{
					TileId:   idCounter,
					Category: _CATEGORY_SUITED,
					Suite:    suit,
					Value:    i + 1,
				})
				idCounter++
			}
		}
	}

	// 2. Generate Honor Tiles (Winds and Dragons)
	// There are 7 honor types. Each tile appears 4 times.
	for _, value := range _ALL_DRAGONS_VALUES {
		for range 4 {
			tiles = append(tiles, tile{
				TileId:   idCounter,
				Category: _CATEGORY_HONOR,
				Suite:    _SUITE_DRAGON,
				Value:    value,
			})
			idCounter++
		}
	}
	for _, value := range _ALL_WINDS_VALUES {
		for range 4 {
			tiles = append(tiles, tile{
				TileId:   idCounter,
				Category: _CATEGORY_HONOR,
				Suite:    _SUITE_WIND,
				Value:    value,
			})
			idCounter++
		}
	}

	// 3. Generate Bonus Tiles (Flowers and Seasons)
	// There are 8 bonus tiles. Each appears only once.
	for _, value := range _ALL_FLOWERS_VALUES {
		tiles = append(tiles, tile{
			TileId:   idCounter,
			Category: _CATEGORY_BONUS,
			Suite:    _SUITE_FLOWER,
			Value:    value,
		})
		idCounter++
	}
	for _, value := range _ALL_SEASONS_VALUES {
		tiles = append(tiles, tile{
			TileId:   idCounter,
			Category: _CATEGORY_BONUS,
			Suite:    _SUITE_SEASON,
			Value:    value,
		})
		idCounter++
	}

	return tiles
}
